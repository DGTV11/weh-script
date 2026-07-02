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
