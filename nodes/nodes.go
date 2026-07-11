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

type StringNode struct {
	BaseNode
	Tok tokens.Token
}

func NewStringNode(tok tokens.Token) StringNode {
	return StringNode{
		Tok:      tok,
		BaseNode: BaseNode{PosRange: position.PositionRange{Start: tok.PosRange.Start, End: tok.PosRange.End}},
	}
}
func (n StringNode) String() string {
	return fmt.Sprintf("%v", n.Tok)
}

type ListNode struct {
	BaseNode
	ElementNodes []Node
}

func NewListNode(elementNodes []Node, posStart *position.Position, posEnd *position.Position) ListNode {
	return ListNode{
		ElementNodes: elementNodes,
		BaseNode:     BaseNode{position.PositionRange{Start: posStart, End: posEnd}},
	}
}
func (n ListNode) String() string {
	return fmt.Sprintf("%v", n.ElementNodes)
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

type VariableDeleteNode struct {
	BaseNode
	VarNameTok tokens.Token
}

func NewVariableDeleteNode(varNameTok tokens.Token) VariableDeleteNode {
	return VariableDeleteNode{
		VarNameTok: varNameTok,
		BaseNode:   BaseNode{PosRange: position.PositionRange{Start: varNameTok.PosRange.Start, End: varNameTok.PosRange.End}},
	}
}
func (n VariableDeleteNode) String() string {
	return fmt.Sprintf("(DELETE %v)", n.VarNameTok)
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
	Cond             Node
	Expr             Node
	ShouldReturnNull bool
}
type ElseCase struct {
	Expr             Node
	ShouldReturnNull bool
}
type IfNode struct {
	BaseNode
	Cases    []IfCase
	ElseCase *ElseCase
}

func NewIfNode(draftIfNode IfNode) IfNode { //this is what a statically typed lang does to people (or just skill issue idk man)
	var lastExpr Node

	cases := draftIfNode.Cases
	elseCase := draftIfNode.ElseCase

	if elseCase == nil {
		// lastNode = cases[len(cases)-1].Cond
		lastExpr = cases[len(cases)-1].Expr
	} else {
		lastExpr = elseCase.Expr
	}

	return IfNode{
		Cases:    cases,
		ElseCase: elseCase,
		BaseNode: BaseNode{PosRange: position.PositionRange{Start: cases[0].Cond.GetPosRange().Start, End: lastExpr.GetPosRange().End}},
	}
}
func (c IfCase) String() string {
	return fmt.Sprintf("(%v : %v ? %t)", c.Cond, c.Expr, c.ShouldReturnNull)
}
func (c ElseCase) String() string {
	return fmt.Sprintf("(%v ? %t)", c.Expr, c.ShouldReturnNull)
}
func (n IfNode) String() string {
	if n.ElseCase == nil {
		return fmt.Sprintf("(IF CASES %v)", n.Cases)
	}
	return fmt.Sprintf("(IF CASES %v ELSE %v)", n.Cases, *n.ElseCase)
}

type ForNode struct {
	BaseNode
	VarNameTok       tokens.Token
	StartValueNode   Node
	StopValueNode    Node
	StepValueNode    Node
	BodyNode         Node
	ShouldReturnNull bool
}

func NewForNode(varNameTok tokens.Token, startValueNode Node, stopValueNode Node, stepValueNode Node, bodyNode Node, shouldReturnNull bool) ForNode {
	return ForNode{
		VarNameTok:       varNameTok,
		StartValueNode:   startValueNode,
		StopValueNode:    stopValueNode,
		StepValueNode:    stepValueNode,
		BodyNode:         bodyNode,
		ShouldReturnNull: shouldReturnNull,
		BaseNode:         BaseNode{PosRange: position.PositionRange{Start: varNameTok.PosRange.Start, End: bodyNode.GetPosRange().End}},
	}
}

func (n ForNode) String() string {
	return fmt.Sprintf("(FOR %v=%v TO %v STEP %v THEN %v ? %t)", n.VarNameTok, n.StartValueNode, n.StopValueNode, n.StepValueNode, n.BodyNode, n.ShouldReturnNull)
}

type WhileNode struct {
	BaseNode
	CondNode         Node
	BodyNode         Node
	ShouldReturnNull bool
}

func NewWhileNode(condNode Node, bodyNode Node, shouldReturnNull bool) WhileNode {
	return WhileNode{
		CondNode:         condNode,
		BodyNode:         bodyNode,
		ShouldReturnNull: shouldReturnNull,
		BaseNode:         BaseNode{PosRange: position.PositionRange{Start: condNode.GetPosRange().Start, End: bodyNode.GetPosRange().End}},
	}
}

func (n WhileNode) String() string {
	return fmt.Sprintf("(WHILE %v THEN %v)", n.CondNode, n.BodyNode)
}

type FuncDefNode struct {
	BaseNode
	VarNameTok       *tokens.Token
	ArgNameToks      []tokens.Token
	BodyNode         Node
	ShouldReturnNull bool
}

func NewFuncDefNode(varNameTok *tokens.Token, argNameToks []tokens.Token, bodyNode Node, shouldReturnNull bool) FuncDefNode {
	var posStart *position.Position

	if varNameTok != nil {
		posStart = varNameTok.PosRange.Start
	} else if len(argNameToks) > 0 {
		posStart = argNameToks[0].PosRange.Start
	} else {
		posStart = bodyNode.GetPosRange().Start
	}

	return FuncDefNode{
		VarNameTok:       varNameTok,
		ArgNameToks:      argNameToks,
		BodyNode:         bodyNode,
		ShouldReturnNull: shouldReturnNull,
		BaseNode:         BaseNode{PosRange: position.PositionRange{Start: posStart, End: bodyNode.GetPosRange().End}},
	}
}

func (n FuncDefNode) String() string {
	return fmt.Sprintf("(FUNC %v ARGS %v => %v ? %t)", n.VarNameTok, n.ArgNameToks, n.BodyNode, n.ShouldReturnNull)
}

type CallNode struct {
	BaseNode
	NodeToCall Node
	ArgNodes   []Node
}

func NewCallNode(nodeToCall Node, argNodes []Node) CallNode {
	var lastNode Node

	if len(argNodes) > 0 {
		lastNode = argNodes[len(argNodes)-1]
	} else {
		lastNode = nodeToCall
	}

	return CallNode{
		NodeToCall: nodeToCall,
		ArgNodes:   argNodes,
		BaseNode:   BaseNode{PosRange: position.PositionRange{Start: nodeToCall.GetPosRange().Start, End: lastNode.GetPosRange().End}},
	}
}

func (n CallNode) String() string {
	return fmt.Sprintf("(CALL %v ARGS %v)", n.NodeToCall, n.ArgNodes)
}

type ItemAccessNode struct {
	BaseNode
	NodeToAccess Node
	KeyNode      Node
}

func NewItemAccessNode(nodeToAccess Node, keyNode Node) ItemAccessNode {
	return ItemAccessNode{
		NodeToAccess: nodeToAccess,
		KeyNode:      keyNode,
		BaseNode:     BaseNode{PosRange: position.PositionRange{Start: nodeToAccess.GetPosRange().Start, End: keyNode.GetPosRange().End}},
	}
}

func (n ItemAccessNode) String() string {
	return fmt.Sprintf("(ACCESS %v KEY %v)", n.NodeToAccess, n.KeyNode)
}

type ItemDeleteNode struct {
	BaseNode
	NodeToAccess Node
	KeyNode      Node
}

func NewItemDeleteNode(nodeToAccess Node, keyNode Node) ItemDeleteNode {
	return ItemDeleteNode{
		NodeToAccess: nodeToAccess,
		KeyNode:      keyNode,
		BaseNode:     BaseNode{PosRange: position.PositionRange{Start: nodeToAccess.GetPosRange().Start, End: keyNode.GetPosRange().End}},
	}
}
func (n ItemDeleteNode) String() string {
	return fmt.Sprintf("(DELETE %v KEY %v)", n.NodeToAccess, n.KeyNode)
}
