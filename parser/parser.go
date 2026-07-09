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

//*Main Parser

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

func (p *Parser) ListExpr() *ParseResult {
	res := NewParseResult()
	var elementNodes []nodes.Node
	posStart := p.CurrentToken.PosRange.Start.Copy()

	if p.CurrentToken.Type != tokens.TokenTypeLsquare {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected '['",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	if p.CurrentToken.Type == tokens.TokenTypeRsquare {
		res.RegisterAdvance()
		p.Advance()
	} else {
		elementNodes = append(elementNodes, res.Register(p.Expr()))
		if res.Err != nil {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected ']', 'var', 'del', 'if', 'for', 'while', 'func', integer, float, identifier, '+', '-', '(', '[' or '!'",
				),
			)
		}

		for p.CurrentToken.Type == tokens.TokenTypeComma {
			res.RegisterAdvance()
			p.Advance()

			elementNodes = append(elementNodes, res.Register(p.Expr()))
			if res.Err != nil {
				return res
			}
		}

		if p.CurrentToken.Type != tokens.TokenTypeRsquare {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected ',' or ']'",
				),
			)
		}

		res.RegisterAdvance()
		p.Advance()
	}

	return res.Success(nodes.NewListNode(elementNodes, posStart, p.CurrentToken.PosRange.End.Copy()))
}

func (p *Parser) IfExpr() *ParseResult {
	res := NewParseResult()
	var cases []nodes.IfCase
	var elseCase nodes.Node = nil

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "if") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'if'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()
	condition := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "then") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'then'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()
	expr := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	cases = append(cases, nodes.IfCase{Cond: condition, Expr: expr})

	for p.CurrentToken.Matches(tokens.TokenTypeKeyword, "elif") {
		res.RegisterAdvance()
		p.Advance()
		condition := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}

		if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "then") {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected 'then'",
				),
			)
		}
		res.RegisterAdvance()
		p.Advance()
		expr := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}

		cases = append(cases, nodes.IfCase{Cond: condition, Expr: expr})
	}

	if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "else") {
		res.RegisterAdvance()
		p.Advance()
		elseCase = res.Register(p.Expr())
		if res.Err != nil {
			return res
		}
	}
	return res.Success(nodes.NewIfNode(cases, elseCase))
}

func (p *Parser) ForExpr() *ParseResult {
	res := NewParseResult()

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "for") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'if'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	if p.CurrentToken.Type != tokens.TokenTypeIdentifier {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected identifier",
			),
		)
	}
	varName := *p.CurrentToken
	res.RegisterAdvance()
	p.Advance()

	if p.CurrentToken.Type != tokens.TokenTypeEquals {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected '='",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	startValue := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "to") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'to'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	stopValue := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	var stepValue nodes.Node = nil
	if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "step") {
		res.RegisterAdvance()
		p.Advance()

		stepValue = res.Register(p.Expr())
		if res.Err != nil {
			return res
		}
	}

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "then") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'then'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	body := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	return res.Success(nodes.NewForNode(varName, startValue, stopValue, stepValue, body))
}

func (p *Parser) FuncDef() *ParseResult {
	res := NewParseResult()

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "func") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'func'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	var varNameTok *tokens.Token
	if p.CurrentToken.Type == tokens.TokenTypeIdentifier {
		varNameTok = p.CurrentToken
		res.RegisterAdvance()
		p.Advance()
		if p.CurrentToken.Type != tokens.TokenTypeLparen {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected '('",
				),
			)
		}
	} else {
		varNameTok = nil
		if p.CurrentToken.Type != tokens.TokenTypeLparen {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected identifier or '('",
				),
			)
		}
	}

	res.RegisterAdvance()
	p.Advance()
	var argNameToks []tokens.Token

	if p.CurrentToken.Type == tokens.TokenTypeIdentifier {
		argNameToks = append(argNameToks, *p.CurrentToken)
		res.RegisterAdvance()
		p.Advance()

		for p.CurrentToken.Type == tokens.TokenTypeComma {
			res.RegisterAdvance()
			p.Advance()

			if p.CurrentToken.Type != tokens.TokenTypeIdentifier {
				return res.Failure(
					errors.NewInvalidSyntaxError(
						p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
						"Expected identifier",
					),
				)
			}

			argNameToks = append(argNameToks, *p.CurrentToken)
			res.RegisterAdvance()
			p.Advance()
		}

		if p.CurrentToken.Type != tokens.TokenTypeRparen {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected ',' or ')'",
				),
			)
		}
	} else {
		if p.CurrentToken.Type != tokens.TokenTypeRparen {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected identifier or ')'",
				),
			)
		}
	}

	res.RegisterAdvance()
	p.Advance()

	if p.CurrentToken.Type != tokens.TokenTypeArrow {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected '=>'",
			),
		)
	}

	res.RegisterAdvance()
	p.Advance()

	body := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	return res.Success(
		nodes.NewFuncDefNode(varNameTok, argNameToks, body),
	)
}

func (p *Parser) WhileExpr() *ParseResult {
	res := NewParseResult()

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "while") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'while'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	condition := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "then") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'then'",
			),
		)
	}
	res.RegisterAdvance()
	p.Advance()

	body := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	return res.Success(nodes.NewWhileNode(condition, body))
}

func (p *Parser) Atom() *ParseResult {
	res := NewParseResult()
	tok := *p.CurrentToken

	if tok.Type == tokens.TokenTypeInt || tok.Type == tokens.TokenTypeFloat {
		res.RegisterAdvance()
		p.Advance()
		return res.Success(nodes.NewNumberNode(tok))
	} else if tok.Type == tokens.TokenTypeString {
		res.RegisterAdvance()
		p.Advance()
		return res.Success(nodes.NewStringNode(tok))
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
	} else if tok.Type == tokens.TokenTypeLsquare {
		listExpr := res.Register(p.ListExpr())
		if res.Err != nil {
			return res
		}
		return res.Success(listExpr)
	} else if tok.Matches(tokens.TokenTypeKeyword, "if") {
		ifExpr := res.Register(p.IfExpr())
		if res.Err != nil {
			return res
		}
		return res.Success(ifExpr)
	} else if tok.Matches(tokens.TokenTypeKeyword, "for") {
		forExpr := res.Register(p.ForExpr())
		if res.Err != nil {
			return res
		}
		return res.Success(forExpr)
	} else if tok.Matches(tokens.TokenTypeKeyword, "while") {
		whileExpr := res.Register(p.WhileExpr())
		if res.Err != nil {
			return res
		}
		return res.Success(whileExpr)
	} else if tok.Matches(tokens.TokenTypeKeyword, "func") {
		funcDef := res.Register(p.FuncDef())
		if res.Err != nil {
			return res
		}
		return res.Success(funcDef)
	}

	return res.Failure(
		errors.NewInvalidSyntaxError(
			tok.PosRange.Start, tok.PosRange.End,
			"Expected integer, float, identifier, '+', '-', '(', '[', 'if', 'for', 'while', or 'func'",
		),
	)
}

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Call, []tokens.TokenType{tokens.TokenTypePow}, p.Factor)
}

func (p *Parser) Call() *ParseResult {
	res := NewParseResult()
	accessNode := res.Register(p.Atom())
	if res.Err != nil {
		return res
	}

	for p.CurrentToken.Type == tokens.TokenTypeLparen || p.CurrentToken.Type == tokens.TokenTypeLsquare {
		switch p.CurrentToken.Type {
		case tokens.TokenTypeLparen:
			res.RegisterAdvance()
			p.Advance()
			var argNodes []nodes.Node

			if p.CurrentToken.Type == tokens.TokenTypeRparen {
				res.RegisterAdvance()
				p.Advance()
			} else {
				argNodes = append(argNodes, res.Register(p.Expr()))
				if res.Err != nil {
					return res.Failure(
						errors.NewInvalidSyntaxError(
							p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
							"Expected ')', 'var', 'del', 'if', 'for', 'while', 'func', integer, float, identifier, '+', '-', '(', '[' or '!'",
						),
					)
				}

				for p.CurrentToken.Type == tokens.TokenTypeComma {
					res.RegisterAdvance()
					p.Advance()

					argNodes = append(argNodes, res.Register(p.Expr()))
					if res.Err != nil {
						return res
					}
				}

				if p.CurrentToken.Type != tokens.TokenTypeRparen {
					return res.Failure(
						errors.NewInvalidSyntaxError(
							p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
							"Expected ',' or ')'",
						),
					)
				}

				res.RegisterAdvance()
				p.Advance()
			}
			accessNode = nodes.NewCallNode(accessNode, argNodes)
		case tokens.TokenTypeLsquare:
			res.RegisterAdvance()
			p.Advance()

			key := res.Register(p.Expr())
			if res.Err != nil {
				return res
			}

			// res.RegisterAdvance()
			// p.Advance()

			if p.CurrentToken.Type != tokens.TokenTypeRsquare {
				return res.Failure(
					errors.NewInvalidSyntaxError(
						p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
						"Expected ']'",
					),
				)
			}

			res.RegisterAdvance()
			p.Advance()

			accessNode = nodes.NewItemAccessNode(accessNode, key)
		}
	}
	return res.Success(accessNode)
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
				"Expected integer, float, identifier, '+', '-', '(', '[' or '!'",
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
	} else if tok.Matches(tokens.TokenTypeKeyword, "del") {
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

		res.RegisterAdvance()
		p.Advance()

		if p.CurrentToken.Type != tokens.TokenTypeLsquare {
			return res.Success(nodes.NewVariableDeleteNode(tok))
		}

		var delNode nodes.Node = nodes.NewVariableAccessNode(tok)

		for p.CurrentToken.Type == tokens.TokenTypeLsquare {
			switch p.CurrentToken.Type {
			case tokens.TokenTypeLsquare:
				res.RegisterAdvance()
				p.Advance()

				key := res.Register(p.Expr())
				if res.Err != nil {
					return res
				}

				// res.RegisterAdvance()
				// p.Advance()

				if p.CurrentToken.Type != tokens.TokenTypeRsquare {
					return res.Failure(
						errors.NewInvalidSyntaxError(
							p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
							"Expected ']'",
						),
					)
				}

				res.RegisterAdvance()
				p.Advance()

				delNode = nodes.NewItemAccessNode(delNode, key)
			}
		}
		delNode = nodes.NewItemDeleteNode(delNode.(nodes.ItemAccessNode).NodeToAccess, delNode.(nodes.ItemAccessNode).KeyNode)
		res.RegisterAdvance()
		p.Advance()

		return res.Success(delNode)
	}

	node := res.Register(p.BinOp(p.CompExpr, []tokens.TokenType{tokens.TokenTypeLAnd, tokens.TokenTypeLOr}, nil))
	if res.Err != nil {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				tok.PosRange.Start, tok.PosRange.End,
				"Expected 'var', 'del', 'if', 'for', 'while', 'func', integer, float, identifier, '+', '-', '(', '[' or '!'",
			),
		)
	}
	return res.Success(node)
}

//*BinOp helpers

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
