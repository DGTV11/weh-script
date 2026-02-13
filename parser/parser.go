package parser

import (
	"slices"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/tokens"
)

type ParseResult struct {
	Err  *errors.Error
	Node nodes.Node
}

func NewParseResult() *ParseResult {
	return &ParseResult{Err: nil, Node: nil}
}

func (pr *ParseResult) Register(res *ParseResult) nodes.Node {
	if res.Err != nil {
		pr.Err = res.Err
	}
	return res.Node
}

func (pr *ParseResult) Success(node nodes.Node) *ParseResult {
	pr.Node = node
	return pr
}

func (pr *ParseResult) Failure(err *errors.Error) *ParseResult {
	pr.Err = err
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
				p.CurrentToken.PosStart, p.CurrentToken.PosEnd,
				"Expected '+', '-', '*', or '/'",
			),
		)
	}
	return res
}

func (p *Parser) Factor() *ParseResult {
	res := NewParseResult()

	tok := *p.CurrentToken

	if tok.Type == tokens.TokenTypePlus || tok.Type == tokens.TokenTypeMinus {
		//TODO: res.Register(p.Advance())
		p.Advance()
		factor := res.Register(p.Factor())
		if res.Err != nil {
			return res
		}
		return res.Success(nodes.UnaryOpNode{OpTok: tok, NodeValue: factor})
	} else if tok.Type == tokens.TokenTypeInt || tok.Type == tokens.TokenTypeFloat {
		//TODO: res.Register(p.Advance())
		p.Advance()
		return res.Success(nodes.NumberNode{Tok: tok})
	} else if tok.Type == tokens.TokenTypeLparen {
		//TODO: res.Register(p.Advance())
		p.Advance()
		expr := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}
		if p.CurrentToken.Type == tokens.TokenTypeRparen {
			//TODO: res.Register(p.Advance())
			p.Advance()
			return res.Success(expr)
		}
		return res.Failure(
			errors.NewInvalidSyntaxError(
				tok.PosStart, tok.PosEnd,
				"Expected ')'",
			),
		)
	}

	return res.Failure(
		errors.NewInvalidSyntaxError(
			tok.PosStart, tok.PosEnd,
			"Expected int or float",
		),
	)
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []tokens.TokenType{tokens.TokenTypeMul, tokens.TokenTypeDiv})
}

func (p *Parser) Expr() *ParseResult {
	return p.BinOp(p.Term, []tokens.TokenType{tokens.TokenTypePlus, tokens.TokenTypeMinus})
}

func (p *Parser) BinOp(function func() *ParseResult, ops []tokens.TokenType) *ParseResult {
	res := NewParseResult()
	var left nodes.Node = nil

	initialLeft := res.Register(function())
	if res.Err != nil {
		return res
	}

	// for p.CurrentToken.Type == tokens.TokenTypeMul || p.CurrentToken.Type == tokens.TokenTypeDiv {
	for slices.Contains(ops, p.CurrentToken.Type) {
		opTok := *p.CurrentToken
		p.Advance() //TODO: res.Register(p.Advance)

		right := res.Register(function())
		if res.Err != nil {
			return res
		}

		if left == nil {
			left = nodes.BinOpNode{LeftNode: initialLeft, OpTok: opTok, RightNode: right}
		} else {
			left = nodes.BinOpNode{LeftNode: left, OpTok: opTok, RightNode: right}
		}
	}

	if left == nil {
		return res.Success(initialLeft)
	}

	return res.Success(left)
}
