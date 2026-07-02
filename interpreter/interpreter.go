package interpreter

import (
	"fmt"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/tokens"
	"github.com/DGTV11/weh-script/values"
)

type RuntimeResult struct {
	Value values.BaseValueInterface
	Err   *errors.Error
}

func NewRuntimeResult() *RuntimeResult {
	return &RuntimeResult{Value: nil, Err: nil}
}

func (rr *RuntimeResult) Register(res *RuntimeResult) values.BaseValueInterface {
	if res.Err != nil {
		rr.Err = res.Err
	}
	return res.Value
}

func (rr *RuntimeResult) Success(value values.BaseValueInterface) *RuntimeResult {
	rr.Value = value
	return rr
}

func (rr *RuntimeResult) Failure(err *errors.Error) *RuntimeResult {
	rr.Err = err
	return rr
}

type Context struct {
	DisplayName    string
	Parent         *Context
	ParentEntryPos *position.Position
}

func Visit(node nodes.Node) *RuntimeResult {
	switch n := node.(type) {
	case nodes.NumberNode:
		return VisitNumberNode(node.(nodes.NumberNode))
	case nodes.BinOpNode:
		return VisitBinOpNode(node.(nodes.BinOpNode))
	case nodes.UnaryOpNode:
		return VisitUnaryOpNode(node.(nodes.UnaryOpNode))
	default:
		posRange := node.GetPosRange()
		return NewRuntimeResult().Failure(errors.NotImplementedError(posRange.Start, posRange.End, fmt.Sprintf("No Visit function defined for node type %T", n)))
	}
}

func VisitNumberNode(node nodes.NumberNode) *RuntimeResult {
	res := NewRuntimeResult()
	var number values.BaseValueInterface = nil

	switch node.Tok.Type {
	case tokens.TokenTypeInt:
		number = &values.Integer{Value: node.Tok.Value.(int64)}
	case tokens.TokenTypeFloat:
		number = &values.Float{Value: node.Tok.Value.(float64)}
	default:
		return res.Failure(errors.NotImplementedError(node.Tok.PosRange.Start, node.Tok.PosRange.End, fmt.Sprintf("NumberNode not implemented for token type %s", tokens.TokenTypeName[node.Tok.Type])))
	}

	number.SetValuePos(node.GetPosRange())
	return res.Success(number)
}

func VisitBinOpNode(node nodes.BinOpNode) *RuntimeResult {
	res := NewRuntimeResult()

	left := res.Register(Visit(node.LeftNode))
	if res.Err != nil {
		return res
	}
	right := res.Register(Visit(node.RightNode))
	if res.Err != nil {
		return res
	}

	var result values.BaseValueInterface = nil
	var error *errors.Error = nil

	switch node.OpTok.Type {
	case tokens.TokenTypePlus:
		result, error = left.Add(right)
	case tokens.TokenTypeMinus:
		result, error = left.Sub(right)
	case tokens.TokenTypeMul:
		result, error = left.Mul(right)
	case tokens.TokenTypeDiv:
		result, error = left.Div(right)
	}

	if error != nil {
		return res.Failure(error)
	}
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitUnaryOpNode(node nodes.UnaryOpNode) *RuntimeResult {
	res := NewRuntimeResult()

	number := res.Register(Visit(node.NodeValue))
	if res.Err != nil {
		return res
	}

	var error *errors.Error = nil

	switch node.OpTok.Type {
	case tokens.TokenTypeMinus:
		number, error = number.Mul(&values.Integer{Value: -1})
	}

	if error != nil {
		return res.Failure(error)
	}
	number.SetValuePos(number.GetPosRange())
	return res.Success(number)
}
