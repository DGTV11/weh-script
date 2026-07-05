package nodes

import (
	"fmt"

	"github.com/DGTV11/weh-script/position"
	"github.com/DGTV11/weh-script/tokens"
)

//*Base Node

type Node interface {
	GetPosRange() position.PositionRange
	String() string
}

type BaseNode struct {
	PosRange position.PositionRange
}

func (b BaseNode) GetPosRange() position.PositionRange {
	return b.PosRange
}

// *Node Definitions
type NumberNode struct {
	BaseNode
	Tok tokens.Token
}

func NewNumberNode(tok tokens.Token) NumberNode {
	return NumberNode{
		Tok:      tok,
		BaseNode: BaseNode{PosRange: position.PositionRange{Start: tok.PosRange.Start, End: tok.PosRange.End}},
	}
}
func (n NumberNode) String() string {
	return fmt.Sprintf("%v", n.Tok)
}

type VariableAccessNode struct {
	BaseNode
	VarNameTok tokens.Token
}

func NewVariableAccessNode(varNameTok tokens.Token) VariableAccessNode {
	return VariableAccessNode{
		VarNameTok: varNameTok,
		BaseNode:   BaseNode{PosRange: position.PositionRange{Start: varNameTok.PosRange.Start, End: varNameTok.PosRange.End}},
	}
}
func (n VariableAccessNode) String() string {
	return fmt.Sprintf("(ACCESS %v)", n.VarNameTok)
}

type VariableAssignNode struct {
	BaseNode
	VarNameTok tokens.Token
	ValueNode  Node
}

func NewVariableAssignNode(varNameTok tokens.Token, valueNode Node) VariableAssignNode {
	return VariableAssignNode{
		VarNameTok: varNameTok,
		ValueNode:  valueNode,
		BaseNode:   BaseNode{PosRange: position.PositionRange{Start: varNameTok.PosRange.Start, End: varNameTok.PosRange.End}},
	}
}
func (n VariableAssignNode) String() string {
	return fmt.Sprintf("(ASSIGN %v = %v)", n.VarNameTok, n.ValueNode)
}

type BinOpNode struct {
	BaseNode
	LeftNode  Node
	OpTok     tokens.Token
	RightNode Node
}

func NewBinOpNode(leftNode Node, opTok tokens.Token, rightNode Node) BinOpNode {
	leftNodePosRange := leftNode.GetPosRange()
	rightNodePosRange := rightNode.GetPosRange()

	return BinOpNode{
		LeftNode:  leftNode,
		OpTok:     opTok,
		RightNode: rightNode,
		BaseNode:  BaseNode{PosRange: position.PositionRange{Start: leftNodePosRange.Start, End: rightNodePosRange.End}},
	}
}
func (n BinOpNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", n.LeftNode, n.OpTok, n.RightNode)
}

type UnaryOpNode struct {
	BaseNode
	OpTok     tokens.Token
	NodeValue Node
}

func NewUnaryOpNode(opTok tokens.Token, nodeValue Node) UnaryOpNode {
	nodeValuePosRange := nodeValue.GetPosRange()
	return UnaryOpNode{
		OpTok:     opTok,
		NodeValue: nodeValue,
		BaseNode:  BaseNode{PosRange: position.PositionRange{Start: opTok.PosRange.Start, End: nodeValuePosRange.End}},
	}
}

func (n UnaryOpNode) String() string {
	return fmt.Sprintf("(%v, %v)", n.OpTok, n.NodeValue)
}

type IfCase struct {
	Cond Node
	Expr Node
}
type IfNode struct {
	BaseNode
	Cases    []IfCase
	ElseCase Node
}

func NewIfNode(cases []IfCase, elseCase Node) IfNode {
	var lastNode Node

	if elseCase == nil {
		lastNode = cases[len(cases)-1].Cond
	} else {
		lastNode = elseCase
	}

	return IfNode{
		Cases:    cases,
		ElseCase: elseCase,
		BaseNode: BaseNode{PosRange: position.PositionRange{Start: cases[0].Cond.GetPosRange().Start, End: lastNode.GetPosRange().End}},
	}
}

func (n IfNode) String() string {
	return fmt.Sprintf("(IF NODE)") //TODO
}
