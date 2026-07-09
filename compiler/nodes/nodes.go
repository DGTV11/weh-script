package nodes

import (
	"fmt"

	"github.com/DGTV11/weh-script/compiler/position"
	"github.com/DGTV11/weh-script/compiler/tokens"
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
	return fmt.Sprintf("(Delete %v)", n.VarNameTok)
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
		// lastNode = cases[len(cases)-1].Cond
		lastNode = cases[len(cases)-1].Expr
	} else {
		lastNode = elseCase
	}

	return IfNode{
		Cases:    cases,
		ElseCase: elseCase,
		BaseNode: BaseNode{PosRange: position.PositionRange{Start: cases[0].Cond.GetPosRange().Start, End: lastNode.GetPosRange().End}},
	}
}
func (c IfCase) String() string {
	return fmt.Sprintf("(%v : %v)", c.Cond, c.Expr)
}
func (n IfNode) String() string {
	return fmt.Sprintf("(IF CASES %v ELSE %v)", n.Cases, n.ElseCase)
}

type ForNode struct {
	BaseNode
	VarNameTok     tokens.Token
	StartValueNode Node
	StopValueNode  Node
	StepValueNode  Node
	BodyNode       Node
}

func NewForNode(varNameTok tokens.Token, startValueNode Node, stopValueNode Node, stepValueNode Node, bodyNode Node) ForNode {
	return ForNode{
		VarNameTok:     varNameTok,
		StartValueNode: startValueNode,
		StopValueNode:  stopValueNode,
		StepValueNode:  stepValueNode,
		BodyNode:       bodyNode,
		BaseNode:       BaseNode{PosRange: position.PositionRange{Start: varNameTok.PosRange.Start, End: bodyNode.GetPosRange().End}},
	}
}

func (n ForNode) String() string {
	return fmt.Sprintf("(FOR %v=%v TO %v STEP %v THEN %v)", n.VarNameTok, n.StartValueNode, n.StopValueNode, n.StepValueNode, n.BodyNode)
}

type WhileNode struct {
	BaseNode
	CondNode Node
	BodyNode Node
}

func NewWhileNode(condNode Node, bodyNode Node) WhileNode {
	return WhileNode{
		CondNode: condNode,
		BodyNode: bodyNode,
		BaseNode: BaseNode{PosRange: position.PositionRange{Start: condNode.GetPosRange().Start, End: bodyNode.GetPosRange().End}},
	}
}

func (n WhileNode) String() string {
	return fmt.Sprintf("(WHILE %v THEN %v)", n.CondNode, n.BodyNode)
}

type FuncDefNode struct {
	BaseNode
	VarNameTok  *tokens.Token
	ArgNameToks []tokens.Token
	BodyNode    Node
}

func NewFuncDefNode(varNameTok *tokens.Token, argNameToks []tokens.Token, bodyNode Node) FuncDefNode {
	var posStart *position.Position

	if varNameTok != nil {
		posStart = varNameTok.PosRange.Start
	} else if len(argNameToks) > 0 {
		posStart = argNameToks[0].PosRange.Start
	} else {
		posStart = bodyNode.GetPosRange().Start
	}

	return FuncDefNode{
		VarNameTok:  varNameTok,
		ArgNameToks: argNameToks,
		BodyNode:    bodyNode,
		BaseNode:    BaseNode{PosRange: position.PositionRange{Start: posStart, End: bodyNode.GetPosRange().End}},
	}
}

func (n FuncDefNode) String() string {
	return fmt.Sprintf("(FUNC %v ARGS %v => %v)", n.VarNameTok, n.ArgNameToks, n.BodyNode)
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
