package parser

import (
	"fmt"
	"slices"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/position"
	"github.com/DGTV11/weh-script/tokens"
)

type ParseResult struct {
	Node                       nodes.Node
	Err                        *errors.Error
	LastRegisteredAdvanceCount int
	AdvanceCount               int
	ToReverseCount             int
}

func NewParseResult() *ParseResult {
	return &ParseResult{}
}

func (pr *ParseResult) RegisterAdvance() {
	pr.LastRegisteredAdvanceCount = 1
	pr.AdvanceCount++
}

func (pr *ParseResult) Register(res *ParseResult) nodes.Node {
	pr.LastRegisteredAdvanceCount = res.AdvanceCount
	pr.AdvanceCount += res.AdvanceCount
	if res.Err != nil {
		pr.Err = res.Err
	}
	return res.Node
}

func (pr *ParseResult) TryRegister(res *ParseResult) nodes.Node {
	if res.Err != nil {
		pr.ToReverseCount = res.AdvanceCount
		return nil
	}
	return pr.Register(res)
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

func (p *Parser) UpdateCurrentTok() {
	if p.TokenIndex >= 0 && p.TokenIndex < len(p.TokenList) {
		p.CurrentToken = &p.TokenList[p.TokenIndex]
	}
}

func (p *Parser) Advance() *tokens.Token {
	p.TokenIndex++
	p.UpdateCurrentTok()
	// fmt.Println(tokens.TokenTypeNameMap[p.CurrentToken.Type])
	return p.CurrentToken
}

func (p *Parser) Reverse(amount int) *tokens.Token {
	p.TokenIndex -= amount
	p.UpdateCurrentTok()
	return p.CurrentToken
}

//*Main Parser

func (p *Parser) Parse() *ParseResult {
	res := p.Statements()

	if res.Err == nil && p.CurrentToken.Type != tokens.TokenTypeEOF {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected '+', '-', '*', '/', '**', '==', '!=', '<', '>', '<=', '>=', '&&' or '||'",
			),
		)
	}
	return res
}

func (p *Parser) Statements() *ParseResult {
	res := NewParseResult()
	var statements []nodes.Node

	posStart := p.CurrentToken.PosRange.Start.Copy()

	for p.CurrentToken.Type == tokens.TokenTypeNewline {
		res.RegisterAdvance()
		p.Advance()
	}

	moreStatements := true

	statement := res.Register(p.Statement())
	if res.Err != nil {
		return res
	}
	// statement := res.TryRegister(p.Statement())
	// if statement == nil {
	// 	p.Reverse(res.ToReverseCount)
	// 	goto finishStatements
	// }
	statements = append(statements, statement)

	for {
		newlineCount := 0
		for p.CurrentToken.Type == tokens.TokenTypeNewline {
			res.RegisterAdvance()
			p.Advance()
			newlineCount += 1
		}
		if newlineCount == 0 {
			moreStatements = false
		}

		if moreStatements == false {
			break
		}
		statement := res.TryRegister(p.Statement())
		if statement == nil {
			p.Reverse(res.ToReverseCount)
			moreStatements = false
			continue
		}
		statements = append(statements, statement)
	}

	// finishStatements:
	return res.Success(nodes.StatementsNode{
		StatementNodes: statements,
		BaseNode:       nodes.BaseNode{position.PositionRange{Start: posStart, End: p.CurrentToken.PosRange.End.Copy()}},
	})
}

func (p *Parser) Statement() *ParseResult {
	res := NewParseResult()
	posStart := p.CurrentToken.PosRange.Start.Copy()
	// fmt.Println("Statement:", p.CurrentToken)
	if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "return") {
		res.RegisterAdvance()
		p.Advance()

		expr := res.TryRegister(p.Expr())
		if expr == nil {
			p.Reverse(res.ToReverseCount)
		}
		return res.Success(nodes.ReturnNode{
			NodeToReturn: expr,
			BaseNode:     nodes.BaseNode{position.PositionRange{Start: posStart, End: p.CurrentToken.PosRange.End.Copy()}},
		})
	} else if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "continue") {
		res.RegisterAdvance()
		p.Advance()
		return res.Success(nodes.ContinueNode{
			BaseNode: nodes.BaseNode{position.PositionRange{Start: posStart, End: p.CurrentToken.PosRange.End.Copy()}},
		})
	} else if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "break") {
		res.RegisterAdvance()
		p.Advance()
		return res.Success(nodes.BreakNode{
			BaseNode: nodes.BaseNode{position.PositionRange{Start: posStart, End: p.CurrentToken.PosRange.End.Copy()}},
		})
	} else if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "import") {
		res.RegisterAdvance()
		p.Advance()

		modulePathTok := *p.CurrentToken
		if modulePathTok.Type != tokens.TokenTypeString {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					modulePathTok.PosRange.Start, modulePathTok.PosRange.End,
					"Expected String literal",
				),
			)
		}
		res.RegisterAdvance()
		p.Advance()

		return res.Success(nodes.NewImportNode(modulePathTok))
	}

	expr := res.Register(p.Expr())
	if res.Err != nil {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'return', 'continue', 'break', 'var', 'del', 'if', 'for', 'while', 'func', integer, float, identifier, '=', '+', '-', '(', '[' or '!'",
			),
		)
	}
	return res.Success(expr)
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
			// return res.Failure(
			// 	errors.NewInvalidSyntaxError(
			// 		tok.PosRange.Start, tok.PosRange.End,
			// 		"Expected '='",
			// 	),
			// )
			return res.Success(nodes.NewVariableAssignNode(varName, nil))
		}

		res.RegisterAdvance()
		p.Advance()

		expr := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}
		return res.Success(nodes.NewVariableAssignNode(varName, expr))
	} else if tok.Matches(tokens.TokenTypeKeyword, "nonlocal") {
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
			// return res.Success(nodes.NewVariableUpdateNode(varName, nil))
		}

		res.RegisterAdvance()
		p.Advance()

		expr := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}
		return res.Success(nodes.NewVariableReassignNode(varName, expr, true))
	} else if tok.Matches(tokens.TokenTypeKeyword, "del") {
		res.RegisterAdvance()
		p.Advance()

		reassignableNode := res.Register(p.Reassignable())
		if res.Err != nil {
			return res
		}

		var delNode nodes.Node
		switch n := reassignableNode.(type) {
		case nodes.VariableAccessNode:
			delNode = nodes.NewVariableDeleteNode(n.VarNameTok)
		case nodes.ItemAccessNode:
			delNode = nodes.NewItemDeleteNode(n.NodeToAccess, n.KeyNode)
		case nodes.MemberAccessNode:
			delNode = nodes.NewMemberDeleteNode(n.NodeToAccess, n.FieldNameTok)
		}

		// res.RegisterAdvance()
		// p.Advance()

		return res.Success(delNode)
	}

	reassignNode := res.TryRegister(p.Reassign())
	if reassignNode != nil {
		return res.Success(reassignNode)
	}
	p.Reverse(res.ToReverseCount)

	// node := res.Register(p.BinOp(p.CompExpr, []tokens.TokenType{tokens.TokenTypeLAnd, tokens.TokenTypeLOr}, nil))
	node := res.Register(p.LOrExpr())
	if res.Err != nil {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				tok.PosRange.Start, tok.PosRange.End,
				"Expected 'var', 'del', 'if', 'for', 'while', 'func', integer, float, identifier, '=', '+', '-', '(', '[' or '!'",
			),
		)
	}
	return res.Success(node)
}

func (p *Parser) Reassign() *ParseResult {
	res := NewParseResult()

	reassignableNode := res.Register(p.Reassignable())
	if res.Err != nil {
		return res
	}

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

	expr := res.Register(p.Expr())
	if res.Err != nil {
		return res
	}

	var reassignNode nodes.Node
	switch n := reassignableNode.(type) {
	case nodes.VariableAccessNode:
		reassignNode = nodes.NewVariableReassignNode(n.VarNameTok, expr, false)
	case nodes.ItemAccessNode:
		reassignNode = nodes.NewItemAssignNode(n.NodeToAccess, n.KeyNode, expr)
	case nodes.MemberAccessNode:
		reassignNode = nodes.NewMemberAssignNode(n.NodeToAccess, n.FieldNameTok, expr)
	}

	// res.RegisterAdvance()
	// p.Advance()

	return res.Success(reassignNode)
}

func (p *Parser) Reassignable() *ParseResult {
	res := NewParseResult()
	tok := *p.CurrentToken
	if p.CurrentToken.Type != tokens.TokenTypeIdentifier {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected identifier",
			),
		)
	}

	res.RegisterAdvance()
	p.Advance()

	var assignableNode nodes.Node = nodes.NewVariableAccessNode(tok)

	for p.CurrentToken.Type == tokens.TokenTypeLsquare || p.CurrentToken.Type == tokens.TokenTypeMemberAccess {
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

			assignableNode = nodes.NewItemAccessNode(assignableNode, key)
		case tokens.TokenTypeMemberAccess:
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

			field := *p.CurrentToken

			res.RegisterAdvance()
			p.Advance()

			assignableNode = nodes.NewMemberAccessNode(assignableNode, field)
		}
	}

	return res.Success(assignableNode)
}

func (p *Parser) LOrExpr() *ParseResult {
	return p.BinOp(p.LAndExpr, []tokens.TokenType{tokens.TokenTypeLOr}, nil)
}

func (p *Parser) LAndExpr() *ParseResult {
	return p.BinOp(p.CompExpr, []tokens.TokenType{tokens.TokenTypeLAnd}, nil)
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

func (p *Parser) ArithExpr() *ParseResult {
	return p.BinOp(p.Term, []tokens.TokenType{tokens.TokenTypePlus, tokens.TokenTypeMinus}, nil)
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []tokens.TokenType{tokens.TokenTypeMul, tokens.TokenTypeDiv}, nil)
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

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Call, []tokens.TokenType{tokens.TokenTypePow}, p.Factor)
}

func (p *Parser) Call() *ParseResult {
	res := NewParseResult()
	accessNode := res.Register(p.Atom())
	if res.Err != nil {
		return res
	}

	for p.CurrentToken.Type == tokens.TokenTypeLparen || p.CurrentToken.Type == tokens.TokenTypeLsquare || p.CurrentToken.Type == tokens.TokenTypeMemberAccess {
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
							"Expected ')', 'var', 'del', 'if', 'for', 'while', 'func', integer, float, identifier, '=', '+', '-', '(', '[' or '!'",
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
		case tokens.TokenTypeMemberAccess:
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

			field := *p.CurrentToken

			res.RegisterAdvance()
			p.Advance()

			accessNode = nodes.NewMemberAccessNode(accessNode, field)
		}
	}
	return res.Success(accessNode)
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
	} else if tok.Type == tokens.TokenTypeChar {
		res.RegisterAdvance()
		p.Advance()
		return res.Success(nodes.NewCharNode(tok))
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
	} else if tok.Matches(tokens.TokenTypeKeyword, "struct") {
		structDef := res.Register(p.StructDef())
		if res.Err != nil {
			return res
		}
		return res.Success(structDef)
	}

	return res.Failure(
		errors.NewInvalidSyntaxError(
			tok.PosRange.Start, tok.PosRange.End,
			"Expected integer, float, identifier, '+', '-', '(', '[', 'if', 'for', 'while', or 'func'",
		),
	)
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
					"Expected ']', 'var', 'del', 'if', 'for', 'while', 'func', integer, float, identifier, '=', '+', '-', '(', '[' or '!'",
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

	return res.Success(nodes.ListNode{
		ElementNodes: elementNodes,
		BaseNode:     nodes.BaseNode{position.PositionRange{Start: posStart, End: p.CurrentToken.PosRange.End.Copy()}},
	})
}

func (p *Parser) IfExprCases(caseKeyword string) *ParseResult {
	res := NewParseResult()
	draftIfNode := nodes.IfNode{}

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, caseKeyword) {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				fmt.Sprintf("Expected '%s'", caseKeyword),
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

	if p.CurrentToken.Type == tokens.TokenTypeNewline {
		res.RegisterAdvance()
		p.Advance()

		statements := res.Register(p.Statements())
		if res.Err != nil {
			return res
		}
		draftIfNode.Cases = append(draftIfNode.Cases, nodes.IfCase{Cond: condition, Expr: statements, ShouldReturnNull: true})

		if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "end") {
			res.RegisterAdvance()
			p.Advance()
		} else {
			newDraftIfNode := res.Register(p.IfExprBOrC())
			if res.Err != nil {
				return res
			}
			draftIfNode.Cases = append(draftIfNode.Cases, newDraftIfNode.(nodes.IfNode).Cases...)
			draftIfNode.ElseCase = newDraftIfNode.(nodes.IfNode).ElseCase
		}
	} else {
		expr := res.Register(p.Statement())
		if res.Err != nil {
			return res
		}
		draftIfNode.Cases = append(draftIfNode.Cases, nodes.IfCase{Cond: condition, Expr: expr, ShouldReturnNull: false})

		newDraftIfNode := res.Register(p.IfExprBOrC())
		if res.Err != nil {
			return res
		}
		draftIfNode.Cases = append(draftIfNode.Cases, newDraftIfNode.(nodes.IfNode).Cases...)
		draftIfNode.ElseCase = newDraftIfNode.(nodes.IfNode).ElseCase
	}
	return res.Success(draftIfNode)
}

func (p *Parser) IfExpr() *ParseResult {
	res := NewParseResult()
	draftIfNode := res.Register(p.IfExprCases("if"))
	if res.Err != nil {
		return res
	}

	return res.Success(nodes.NewIfNode(draftIfNode.(nodes.IfNode)))
}

func (p *Parser) IfExprB() *ParseResult {
	return p.IfExprCases("elif")
}

func (p *Parser) IfExprC() *ParseResult {
	res := NewParseResult()
	draftIfNode := nodes.IfNode{}

	if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "else") {
		res.RegisterAdvance()
		p.Advance()

		if p.CurrentToken.Type == tokens.TokenTypeNewline {
			res.RegisterAdvance()
			p.Advance()

			statements := res.Register(p.Statements())
			if res.Err != nil {
				return res
			}
			draftIfNode.ElseCase = &nodes.ElseCase{Expr: statements, ShouldReturnNull: true}

			if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "end") {
				res.RegisterAdvance()
				p.Advance()
			} else {
				return res.Failure(
					errors.NewInvalidSyntaxError(
						p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
						"Expected 'end'",
					),
				)
			}
		} else {
			expr := res.Register(p.Statement())
			if res.Err != nil {
				return res
			}
			draftIfNode.ElseCase = &nodes.ElseCase{Expr: expr, ShouldReturnNull: false}
		}
	}

	return res.Success(draftIfNode)
}

func (p *Parser) IfExprBOrC() *ParseResult {
	res := NewParseResult()
	draftIfNode := nodes.IfNode{}

	if p.CurrentToken.Matches(tokens.TokenTypeKeyword, "elif") {
		draftIfNode = res.Register(p.IfExprB()).(nodes.IfNode)
	} else {
		draftIfNode = res.Register(p.IfExprC()).(nodes.IfNode)
	}
	if res.Err != nil {
		return res
	}

	return res.Success(draftIfNode)
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

	if p.CurrentToken.Type == tokens.TokenTypeNewline {
		res.RegisterAdvance()
		p.Advance()

		body := res.Register(p.Statements())
		if res.Err != nil {
			return res
		}

		if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "end") {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected 'end'",
				),
			)
		}

		res.RegisterAdvance()
		p.Advance()

		return res.Success(nodes.NewForNode(varName, startValue, stopValue, stepValue, body, true))
	}

	body := res.Register(p.Statement())
	if res.Err != nil {
		return res
	}

	return res.Success(nodes.NewForNode(varName, startValue, stopValue, stepValue, body, false))
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

	if p.CurrentToken.Type == tokens.TokenTypeNewline {
		res.RegisterAdvance()
		p.Advance()

		body := res.Register(p.Statements())
		if res.Err != nil {
			return res
		}

		if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "end") {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected 'end'",
				),
			)
		}

		res.RegisterAdvance()
		p.Advance()

		return res.Success(nodes.NewWhileNode(condition, body, true))
	}

	body := res.Register(p.Statement())
	if res.Err != nil {
		return res
	}

	return res.Success(nodes.NewWhileNode(condition, body, false))
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

	if p.CurrentToken.Type == tokens.TokenTypeArrow {
		res.RegisterAdvance()
		p.Advance()

		body := res.Register(p.Expr())
		if res.Err != nil {
			return res
		}

		return res.Success(
			nodes.NewFuncDefNode(varNameTok, argNameToks, body, true),
		)
	}

	if p.CurrentToken.Type != tokens.TokenTypeNewline {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected '=>' or Newline",
			),
		)
	}

	res.RegisterAdvance()
	p.Advance()

	body := res.Register(p.Statements())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "end") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'end'",
			),
		)
	}

	res.RegisterAdvance()
	p.Advance()

	return res.Success(
		nodes.NewFuncDefNode(varNameTok, argNameToks, body, false),
	)
}

func (p *Parser) StructDef() *ParseResult {
	res := NewParseResult()

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "struct") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'struct'",
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
		if p.CurrentToken.Type != tokens.TokenTypeNewline {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected Newline",
				),
			)
		}
	} else {
		varNameTok = nil
		if p.CurrentToken.Type != tokens.TokenTypeNewline {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected identifier or Newline",
				),
			)
		}
	}

	res.RegisterAdvance()
	p.Advance()

	for p.CurrentToken.Type == tokens.TokenTypeNewline {
		res.RegisterAdvance()
		p.Advance()
	}

	if p.CurrentToken.Type != tokens.TokenTypeIdentifier {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected identifier",
			),
		)
	}

	var fieldNameToks []tokens.Token
	fieldNameToks = append(fieldNameToks, *p.CurrentToken) //TODO: any better way to do this + previous line

	res.RegisterAdvance()
	p.Advance()

	for p.CurrentToken.Type == tokens.TokenTypeNewline {
		res.RegisterAdvance()
		p.Advance()
	}

	if p.CurrentToken.Type != tokens.TokenTypeIdentifier {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected identifier",
			),
		)
	}

	for p.CurrentToken.Type == tokens.TokenTypeIdentifier {
		fieldNameToks = append(fieldNameToks, *p.CurrentToken)
		res.RegisterAdvance()
		p.Advance()
		if p.CurrentToken.Type != tokens.TokenTypeNewline {
			return res.Failure(
				errors.NewInvalidSyntaxError(
					p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
					"Expected Newline",
				),
			)
		}
		for p.CurrentToken.Type == tokens.TokenTypeNewline {
			res.RegisterAdvance()
			p.Advance()
		}
	}

	if !p.CurrentToken.Matches(tokens.TokenTypeKeyword, "end") {
		return res.Failure(
			errors.NewInvalidSyntaxError(
				p.CurrentToken.PosRange.Start, p.CurrentToken.PosRange.End,
				"Expected 'end'",
			),
		)
	}

	res.RegisterAdvance()
	p.Advance()

	return res.Success(
		nodes.NewStructDefNode(varNameTok, fieldNameToks),
	)
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
