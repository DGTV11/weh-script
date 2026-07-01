package nodes

import (
	"fmt"

	// "github.com/DGTV11/weh-script/interpreter"
	"github.com/DGTV11/weh-script/tokens"
)

//*Node definitions

type Node interface {
	String() string
	Eval() any //TODO: RTResult?
}

type NumberNode struct {
	Node
	Tok tokens.Token
}

func (n NumberNode) String() string {
	return fmt.Sprintf("%v", n.Tok)
}

type BinOpNode struct {
	Node
	LeftNode  Node
	OpTok     tokens.Token
	RightNode Node
}

func (n BinOpNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", n.LeftNode, n.OpTok, n.RightNode)
}

type UnaryOpNode struct {
	Node
	OpTok     tokens.Token
	NodeValue Node
}

func (n UnaryOpNode) String() string {
	return fmt.Sprintf("(%v, %v)", n.OpTok, n.NodeValue)
}

//*Interpreter

func (n NumberNode) Eval() any {
	fmt.Printf("Found number node!")
	return nil
}

func (n BinOpNode) Eval() any {
	fmt.Printf("Found bin op node!")
	n.LeftNode.Eval()
	n.RightNode.Eval()
	return nil
}

func (n UnaryOpNode) Eval() any {
	fmt.Printf("Found unary op node!")
	n.NodeValue.Eval()
	return nil
}
