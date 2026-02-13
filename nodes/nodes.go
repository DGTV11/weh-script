package nodes

import (
	"fmt"

	"github.com/DGTV11/weh-script/tokens"
)

type Node interface {
	String() string
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
