package interpreter

import (
	"fmt"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/runtime"
	"github.com/DGTV11/weh-script/tokens"
	"github.com/DGTV11/weh-script/values"
)

// *RuntimeResult
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

// *Main Interpreter
func Visit(node nodes.Node, ctx runtime.Context) *RuntimeResult {
	switch n := node.(type) {
	case nodes.NumberNode:
		return VisitNumberNode(node.(nodes.NumberNode), ctx)
	case nodes.VariableAccessNode:
		return VisitVariableAccessNode(node.(nodes.VariableAccessNode), ctx)
	case nodes.VariableAssignNode:
		return VisitVariableAssignNode(node.(nodes.VariableAssignNode), ctx)
	case nodes.BinOpNode:
		return VisitBinOpNode(node.(nodes.BinOpNode), ctx)
	case nodes.UnaryOpNode:
		return VisitUnaryOpNode(node.(nodes.UnaryOpNode), ctx)
	case nodes.IfNode:
		return VisitIfNode(node.(nodes.IfNode), ctx)
	default:
		posRange := node.GetPosRange()
		return NewRuntimeResult().Failure(errors.NewNotImplementedError(posRange.Start, posRange.End, fmt.Sprintf("No Visit function defined for node type %T", n), ctx))
	}
}

func VisitNumberNode(node nodes.NumberNode, ctx runtime.Context) *RuntimeResult {
	res := NewRuntimeResult()
	var number values.BaseValueInterface = nil

	switch node.Tok.Type {
	case tokens.TokenTypeInt:
		number = &values.Integer{Value: node.Tok.Value.(int64)}
	case tokens.TokenTypeFloat:
		number = &values.Float{Value: node.Tok.Value.(float64)}
	default:
		return res.Failure(errors.NewNotImplementedError(node.Tok.PosRange.Start, node.Tok.PosRange.End, fmt.Sprintf("NumberNode not implemented for token type %s", tokens.TokenTypeName[node.Tok.Type]), ctx))
	}

	number.SetContext(ctx)
	number.SetValuePos(node.GetPosRange())
	return res.Success(number)
}

func VisitVariableAccessNode(node nodes.VariableAccessNode, ctx runtime.Context) *RuntimeResult {
	res := NewRuntimeResult()
	posRange := node.GetPosRange()

	varName := node.VarNameTok.Value.(string)
	rawValue := ctx.SymTable.GetSymbol(varName)

	if rawValue == nil {
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is not defined", varName), ctx))
	}

	value := rawValue.(values.BaseValueInterface)

	value.SetContext(ctx)
	value.SetValuePos(posRange)
	return res.Success(value)
}

func VisitVariableAssignNode(node nodes.VariableAssignNode, ctx runtime.Context) *RuntimeResult {
	res := NewRuntimeResult()
	varName := node.VarNameTok.Value.(string)
	value := res.Register(Visit(node.ValueNode, ctx))
	if res.Err != nil {
		return res
	}

	ctx.SymTable.SetSymbol(varName, value)
	return res.Success(value)
}

func VisitBinOpNode(node nodes.BinOpNode, ctx runtime.Context) *RuntimeResult {
	res := NewRuntimeResult()

	left := res.Register(Visit(node.LeftNode, ctx))
	if res.Err != nil {
		return res
	}
	right := res.Register(Visit(node.RightNode, ctx))
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
	case tokens.TokenTypePow:
		result, error = left.Pow(right)
	case tokens.TokenTypeEE:
		result, error = left.Eq(right)
	case tokens.TokenTypeNE:
		result, error = left.Ne(right)
	case tokens.TokenTypeLT:
		result, error = left.Lt(right)
	case tokens.TokenTypeGT:
		result, error = left.Gt(right)
	case tokens.TokenTypeLTE:
		result, error = left.Lte(right)
	case tokens.TokenTypeGTE:
		result, error = left.Gte(right)
	case tokens.TokenTypeLAnd:
		result, error = left.LAnd(right)
	case tokens.TokenTypeLOr:
		result, error = left.LOr(right)
	}

	if error != nil {
		return res.Failure(error)
	}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitUnaryOpNode(node nodes.UnaryOpNode, ctx runtime.Context) *RuntimeResult {
	res := NewRuntimeResult()

	number := res.Register(Visit(node.NodeValue, ctx))
	if res.Err != nil {
		return res
	}

	var error *errors.Error = nil

	switch node.OpTok.Type {
	case tokens.TokenTypeMinus:
		number, error = number.Mul(&values.Integer{Value: -1})
	case tokens.TokenTypeLNot:
		number, error = number.LNot()
	}

	if error != nil {
		return res.Failure(error)
	}
	number.SetContext(ctx)
	number.SetValuePos(number.GetPosRange())
	return res.Success(number)
}

func VisitIfNode(node nodes.IfNode, ctx runtime.Context) *RuntimeResult {
	res := NewRuntimeResult()

	for i := 0; i < len(node.Cases); i++ {
		conditionValue := res.Register(Visit(node.Cases[i].Cond, ctx))
		if res.Err != nil {
			return res
		}

		if conditionValue.IsTrue() {
			exprValue := res.Register(Visit(node.Cases[i].Expr, ctx))
			if res.Err != nil {
				return res
			}
			return res.Success(exprValue)
		}
	}
	if node.ElseCase != nil {
		elseValue := res.Register(Visit(node.ElseCase, ctx))
		if res.Err != nil {
			return res
		}
		return res.Success(elseValue)
	}

	return res.Success(nil)
}
