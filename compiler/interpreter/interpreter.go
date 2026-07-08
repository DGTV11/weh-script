package interpreter

import (
	"fmt"

	"github.com/DGTV11/weh-script/compiler/environment"
	"github.com/DGTV11/weh-script/compiler/errors"
	"github.com/DGTV11/weh-script/compiler/nodes"
	"github.com/DGTV11/weh-script/compiler/tokens"
	"github.com/DGTV11/weh-script/compiler/values"
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
func Visit(node nodes.Node, ctx environment.Context) *RuntimeResult {
	switch n := node.(type) {
	case nodes.NumberNode:
		return VisitNumberNode(node.(nodes.NumberNode), ctx)
	case nodes.StringNode:
		return VisitStringNode(node.(nodes.StringNode), ctx)
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
	case nodes.ForNode:
		return VisitForNode(node.(nodes.ForNode), ctx)
	case nodes.WhileNode:
		return VisitWhileNode(node.(nodes.WhileNode), ctx)
	case nodes.FuncDefNode:
		return VisitFuncDefNode(node.(nodes.FuncDefNode), ctx)
	case nodes.CallNode:
		return VisitCallNode(node.(nodes.CallNode), ctx)
	default:
		posRange := node.GetPosRange()
		return NewRuntimeResult().Failure(errors.NewNotImplementedError(posRange.Start, posRange.End, fmt.Sprintf("No Visit function defined for node type %T", n), ctx))
	}
}

func VisitNumberNode(node nodes.NumberNode, ctx environment.Context) *RuntimeResult {
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

func VisitStringNode(node nodes.StringNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	str := &values.String{Value: node.Tok.Value.(string)}
	str.SetContext(ctx)
	str.SetValuePos(node.GetPosRange())
	return res.Success(str)
}

func VisitVariableAccessNode(node nodes.VariableAccessNode, ctx environment.Context) *RuntimeResult {
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

func VisitVariableAssignNode(node nodes.VariableAssignNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	varName := node.VarNameTok.Value.(string)
	value := res.Register(Visit(node.ValueNode, ctx))
	if res.Err != nil {
		return res
	}

	ctx.SymTable.SetSymbol(varName, value)
	return res.Success(value)
}

func VisitBinOpNode(node nodes.BinOpNode, ctx environment.Context) *RuntimeResult {
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

func VisitUnaryOpNode(node nodes.UnaryOpNode, ctx environment.Context) *RuntimeResult {
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
	number.SetValuePos(node.GetPosRange())
	return res.Success(number)
}

func VisitIfNode(node nodes.IfNode, ctx environment.Context) *RuntimeResult {
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

	// return res.Success(nil)
	return res.Success(&values.Integer{Value: 0})
}

func VisitForNode(node nodes.ForNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	startValue := res.Register(Visit(node.StartValueNode, ctx))
	if res.Err != nil {
		return res
	}

	stopValue := res.Register(Visit(node.StopValueNode, ctx))
	if res.Err != nil {
		return res
	}

	var stepValue values.BaseValueInterface = &values.Integer{Value: 1}
	if node.StepValueNode != nil {
		stepValue = res.Register(Visit(node.StepValueNode, ctx))
		if res.Err != nil {
			return res
		}
	}

	i := startValue.Copy()

	condRes, err := stepValue.Gte(&values.Integer{Value: 0})
	if err != nil {
		return res.Failure(err)
	}

	var cond func() (values.BaseValueInterface, *errors.Error)
	if condRes.IsTrue() {
		cond = func() (values.BaseValueInterface, *errors.Error) { return i.Lt(stopValue) }
	} else {
		cond = func() (values.BaseValueInterface, *errors.Error) { return i.Gt(stopValue) }
	}

	condRes, err = cond()
	if err != nil {
		return res.Failure(err)
	}

	for condRes.IsTrue() {
		ctx.SymTable.SetSymbol(node.VarNameTok.Value.(string), i)
		i, err = i.Add(stepValue)
		if err != nil {
			return res.Failure(err)
		}

		res.Register(Visit(node.BodyNode, ctx))
		if res.Err != nil {
			return res
		}

		condRes, err = cond()
		if err != nil {
			return res.Failure(err)
		}
	}

	// return res.Success(nil)
	return res.Success(&values.Integer{Value: 0})
}

func VisitWhileNode(node nodes.WhileNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	for {
		condition := res.Register(Visit(node.CondNode, ctx))
		if res.Err != nil {
			return res
		}
		if !condition.IsTrue() {
			break
		}
		res.Register(Visit(node.BodyNode, ctx))
		if res.Err != nil {
			return res
		}
	}

	// return res.Success(nil)
	return res.Success(&values.Integer{Value: 0})
}

func VisitFuncDefNode(node nodes.FuncDefNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	var funcName *string = nil
	if node.VarNameTok != nil {
		funcNameStrVal := node.VarNameTok.Value.(string)
		funcName = &funcNameStrVal
	}

	bodyNode := node.BodyNode
	var argNames []string
	for i := 0; i < len(node.ArgNameToks); i++ {
		argNames = append(argNames, node.ArgNameToks[i].Value.(string))
	}
	funcValue := values.NewFunction(funcName, bodyNode, argNames)
	funcValue.SetContext(ctx)
	funcValue.SetValuePos(node.GetPosRange())

	if funcName != nil {
		ctx.SymTable.SetSymbol(*funcName, funcValue)
	}
	return res.Success(funcValue)
}

func VisitCallNode(node nodes.CallNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	var args []values.BaseValueInterface

	valueToCall := res.Register(Visit(node.NodeToCall, ctx))
	if res.Err != nil {
		return res
	}
	valueToCall = valueToCall.Copy()
	valueToCall.SetValuePos(node.GetPosRange())

	for i := 0; i < len(node.ArgNodes); i++ {
		args = append(args, res.Register(Visit(node.ArgNodes[i], ctx)))
		if res.Err != nil {
			return res
		}
	}
	returnValue := res.Register(ExecuteFunction(valueToCall.(*values.Function), args))
	if res.Err != nil {
		return res
	}
	return res.Success(returnValue)
}

// *Function Calls
func ExecuteFunction(function *values.Function, args []values.BaseValueInterface) *RuntimeResult {
	res := NewRuntimeResult()
	posRange := function.GetPosRange()
	parentCtx := function.GetContext()
	newCtx := environment.Context{DisplayName: function.Name, Parent: &parentCtx, ParentEntryPos: function.GetPosRange().Start, SymTable: parentCtx.SymTable}

	if len(args) > len(function.ArgNames) {
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("%d too many args passed into '%s'", len(args)-len(function.ArgNames), function.Name), parentCtx))
	}
	if len(args) < len(function.ArgNames) {
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("%d too few args passed into '%s'", len(function.ArgNames)-len(args), function.Name), parentCtx))
	}
	for i := 0; i < len(args); i++ {
		argName := function.ArgNames[i]
		argValue := args[i]
		argValue.SetContext(newCtx)
		newCtx.SymTable.SetSymbol(argName, argValue)
	}

	value := res.Register(Visit(function.BodyNode, newCtx))
	if res.Err != nil {
		return res
	}
	return res.Success(value)
}
