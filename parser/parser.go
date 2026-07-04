package parser

import (
	"slices"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/tokens"
)

type ParseResult struct {
	Node         nodes.Node
	Err          *errors.Error
	AdvanceCount int
}

func NewParseResult() *ParseResult {
	return &ParseResult{Node: nil, Err: nil, AdvanceCount: 0}
}

func (pr *ParseResult) Register(res *ParseResult) nodes.Node {
	pr.AdvanceCount += res.AdvanceCount
	if res.Err != nil {
		pr.Err = res.Err
	}
	return res.Node
}

func (pr *ParseResult) RegisterAdvance() {
	pr.AdvanceCount++
}

func (pr *ParseResult) Success(node nodes.Node) *ParseResult {
	pr.Node = node
	return pr
}

func (pr *ParseResult) Failure(err *errors.Error) *ParseResult {
	if pr.Err == nil || pr.AdvanceCount == 0 {
		pr.Err = err
	}
	return pr
}

type Parser struct {
	TokenList    []tokens.Token
	TokenIndex   int
	CurrentToken *tokens.Token
}

func NewParser(tokenList []tokens.Token) *Parser {
	newParser := Parser{
		TokenList:    tokenList,
		TokenIndex:   -1,
		CurrentToken: nil,
	}
	newParser.Advance()
	return &newParser
}

func (p *Parser) Advance() *tokens.Token {
	p.TokenIndex += 1

	if p.TokenIndex < len(p.TokenList) {
		p.CurrentToken = &p.TokenList[p.TokenIndex]
	}
	return p.CurrentToken
}

func (p *Parser) Parse() *ParseResult {
	res := p.Expr()

	if res.Err == nil && p.CurrentToken.Type != tokens.TokenTypeEOF {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected '+', '-', '*', or '/'",
			),
		)
	}
	return res
}

func (p *Parser) Atom() *ParseResult {
	res := NewParseResult()
	tok := *p.CurrentToken

	if tok.Type == tokens.TokenTypeInt || tok.Type == tokens.TokenTypeFloat {
		res.RegisterAdvance()
		p.Advance()
		return res.Success(nodes.NewNumberNode(tok))
	} else if tok.Type == tokens.TokenTypeIdentifier {
		res.RegisterAdvance()
		p.Advance()
		return res.Success(nodes.NewVariableAccessNode(tok))
	} else if tok.Type == tokens.TokenTypeLparen {
		res.RegisterAdvance()
		p.Advance()
		expr := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}
		if p.CurrentToken.Type == tokens.TokenTypeRparen {
			res.RegisterAdvance()
			p.Advance()
			return res.Success(expr)
		}
		return res.Failure(
			errors.NewInvalidSyntaxError(
				tok.PosRange.Start, tok.PosRange.End,
				"Expected ')'",
			),
		)
	}

	return res.Failure(
		errors.NewInvalidSyntaxError(
			tok.PosRange.Start, tok.PosRange.End,
			"Expected integer, float, identifier, '+', '-' or '('",
		),
	)
}

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Atom, []tokens.TokenType{tokens.TokenTypePow}, p.Factor)
}

func (p *Parser) Factor() *ParseResult {
	res := NewParseResult()
	tok := *p.CurrentToken

	if tok.Type == tokens.TokenTypePlus || tok.Type == tokens.TokenTypeMinus {
		res.RegisterAdvance()
		p.Advance()
		factor := res.Register(p.Factor())
		if res.Err != nil {
			return res
		}
		return res.Success(nodes.NewUnaryOpNode(tok, factor))
	}

	return p.Power()
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []tokens.TokenType{tokens.TokenTypeMul, tokens.TokenTypeDiv}, nil)
}

func (p *Parser) ArithExpr() *ParseResult {
	return p.BinOp(p.Term, []tokens.TokenType{tokens.TokenTypePlus, tokens.TokenTypeMinus}, nil)
}

func (p *Parser) CompExpr() *ParseResult {
	res := NewParseResult()
	tok := *p.CurrentToken

	if tok.Type == tokens.TokenTypeLNot {
		res.RegisterAdvance()
		p.Advance()

		node := res.Register(p.CompExpr())
		if res.Err != nil {
			return res
		}
		return res.Success(nodes.NewUnaryOpNode(tok, node))
	}

	node := res.Register(p.BinOp(p.ArithExpr, []tokens.TokenType{tokens.TokenTypeEE, tokens.TokenTypeNE, tokens.TokenTypeLT, tokens.TokenTypeGT, tokens.TokenTypeLTE, tokens.TokenTypeGTE}, nil))
	if res.Err != nil {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				tok.PosRange.Start, tok.PosRange.End,
				"Expected integer, float, identifier, '+', '-', '(', '!'",
			),
		)
	}

	return res.Success(node)
}

func (p *Parser) Expr() *ParseResult {
	res := NewParseResult()
	tok := *p.CurrentToken

	if tok.Matches(tokens.TokenTypeKeyword, "var") {
		res.RegisterAdvance()
		p.Advance()

		tok = *p.CurrentToken
		if tok.Type != tokens.TokenTypeIdentifier {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					tok.PosRange.Start, tok.PosRange.End,
					"Expected identifier",
				),
			)
		}

		varName := tok
		res.RegisterAdvance()
		p.Advance()

		tok = *p.CurrentToken
		if tok.Type != tokens.TokenTypeEquals {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					tok.PosRange.Start, tok.PosRange.End,
					"Expected '='",
				),
			)
		}

		res.RegisterAdvance()
		p.Advance()

		expr := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}
		return res.Success(nodes.NewVariableAssignNode(varName, expr))
	}

	node := res.Register(p.BinOp(p.CompExpr, []tokens.TokenType{tokens.TokenTypeLAnd, tokens.TokenTypeLOr}, nil))
	if res.Err != nil {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				tok.PosRange.Start, tok.PosRange.End,
				"Expected 'var', integer, float, identifier, '+', '-' or '('",
			),
		)
	}
	return res.Success(node)
}

func (p *Parser) BinOp(functionL func() *ParseResult, ops []tokens.TokenType, functionR func() *ParseResult) *ParseResult {
	if functionR == nil {
		functionR = functionL
	}

	res := NewParseResult()
	var left nodes.Node = nil

	initialLeft := res.Register(functionL())
	if res.Err != nil {
		return res
	}

	for slices.Contains(ops, p.CurrentToken.Type) {
		opTok := *p.CurrentToken
		res.RegisterAdvance()
		p.Advance()

		right := res.Register(functionR())
		if res.Err != nil {
			return res
		}

		if left == nil {
			left = nodes.NewBinOpNode(initialLeft, opTok, right)
		} else {
			left = nodes.NewBinOpNode(left, opTok, right)
		}
	}

	if left == nil {
		return res.Success(initialLeft)
	}

	return res.Success(left)
}

func (p *Parser) BinOpWithTokTVs(functionL func() *ParseResult, ops []tokens.TokenTV, functionR func() *ParseResult) *ParseResult {
	if functionR == nil {
		functionR = functionL
	}

	res := NewParseResult()
	var left nodes.Node = nil

	initialLeft := res.Register(functionL())
	if res.Err != nil {
		return res
	}

	for slices.Contains(ops, tokens.TokenTV{Type: p.CurrentToken.Type, Value: p.CurrentToken.Value}) {
		opTok := *p.CurrentToken
		res.RegisterAdvance()
		p.Advance()

		right := res.Register(functionR())
		if res.Err != nil {
			return res
		}

		if left == nil {
			left = nodes.NewBinOpNode(initialLeft, opTok, right)
		} else {
			left = nodes.NewBinOpNode(left, opTok, right)
		}
	}

	if left == nil {
		return res.Success(initialLeft)
	}

	return res.Success(left)
}
