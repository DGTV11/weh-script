package nodes

import (
	"fmt"

	// "github.com/DGTV11/weh-script/interpreter"
	"github.com/DGTV11/weh-script/position"
	"github.com/DGTV11/weh-script/tokens"
	"github.com/DGTV11/weh-script/values"
)

//*Node definitions

type Node interface {
	GetPosRange() position.PositionRange
	String() string
	Eval() values.BaseValueInterface //TODO: RTResult?
}

type BaseNode struct {
	PosRange position.PositionRange
}

func (b BaseNode) GetPosRange() position.PositionRange {
	return b.PosRange
}

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

//*Interpreter

func (n NumberNode) Eval() values.BaseValueInterface {
	var result values.BaseValueInterface = nil

	switch n.Tok.Type {
	case tokens.TokenTypeInt:
		result = &values.Integer{Value: n.Tok.Value.(int64)}
	case tokens.TokenTypeFloat:
		result = &values.Float{Value: n.Tok.Value.(float64)}
	}

	result.SetValuePos(result.GetPosRange())
	return result
}

func (n BinOpNode) Eval() values.BaseValueInterface {
	left := n.LeftNode.Eval()
	right := n.RightNode.Eval()

	var result values.BaseValueInterface = nil

	switch n.OpTok.Type {
	case tokens.TokenTypePlus:
		result = left.Add(right)
	case tokens.TokenTypeMinus:
		result = left.Sub(right)
	case tokens.TokenTypeMul:
		result = left.Mul(right)
	case tokens.TokenTypeDiv:
		result = left.Div(right)
	}

	result.SetValuePos(result.GetPosRange())
	return result
}

func (n UnaryOpNode) Eval() values.BaseValueInterface {
	number := n.NodeValue.Eval()

	switch n.OpTok.Type {
	case tokens.TokenTypeMinus:
		number = number.Mul(&values.Integer{Value: -1})
	}

	number.SetValuePos(number.GetPosRange())
	return number
}
