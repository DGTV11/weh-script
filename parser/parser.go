package parser

import (
	"slices"

	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/tokens"
)

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

func (p *Parser) Advance() {
	p.TokenIndex += 1

	if p.TokenIndex < len(p.TokenList) {
		p.CurrentToken = &p.TokenList[p.TokenIndex]
	}
}

func (p *Parser) Parse() nodes.Node {
	res := p.Expr()
	return res
}

func (p *Parser) Factor() nodes.Node {
	tok := *p.CurrentToken

	if tok.Type == tokens.TokenTypeInt || tok.Type == tokens.TokenTypeFloat {
		p.Advance()
		return nodes.NumberNode{Tok: tok}
	}

	return nil //TODO: error handling
}

func (p *Parser) Term() nodes.Node {
	return p.BinOp(p.Factor, []tokens.TokenType{tokens.TokenTypeMul, tokens.TokenTypeDiv})
}

func (p *Parser) Expr() nodes.Node {
	return p.BinOp(p.Term, []tokens.TokenType{tokens.TokenTypePlus, tokens.TokenTypeMinus})
}

func (p *Parser) BinOp(function func() nodes.Node, ops []tokens.TokenType) nodes.Node {
	var left nodes.Node = nil

	initialLeft := function()
	if initialLeft == nil {
		//TODO: error handling
	}

	// for p.CurrentToken.Type == tokens.TokenTypeMul || p.CurrentToken.Type == tokens.TokenTypeDiv {
	for slices.Contains(ops, p.CurrentToken.Type) {
		opTok := *p.CurrentToken
		p.Advance()

		right := function()
		if right == nil {
			//TODO: error handling
		}

		if left == nil {
			left = nodes.BinOpNode{LeftNode: initialLeft, OpTok: opTok, RightNode: right}
		} else {
			left = nodes.BinOpNode{LeftNode: left, OpTok: opTok, RightNode: right}
		}
	}

	if left == nil {
		return initialLeft
	}

	return left
}
