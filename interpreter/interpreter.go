package interpreter

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/inancgumus/screen"

	"github.com/DGTV11/weh-script/environment"
	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/tokens"
	"github.com/DGTV11/weh-script/values"
)

// *RuntimeResult
type RuntimeResult struct {
	Value              values.BaseValueInterface
	Err                *errors.Error
	FuncReturnValue    values.BaseValueInterface
	LoopShouldContinue bool
	LoopShouldBreak    bool
}

func NewRuntimeResult() *RuntimeResult {
	return &RuntimeResult{}
}

func (rr *RuntimeResult) Reset() {
	rr.Value = nil
	rr.Err = nil
	rr.FuncReturnValue = nil
	rr.LoopShouldContinue = false
	rr.LoopShouldBreak = false
}

func (rr *RuntimeResult) Register(res *RuntimeResult) values.BaseValueInterface {
	rr.Err = res.Err
	rr.FuncReturnValue = res.FuncReturnValue
	rr.LoopShouldContinue = res.LoopShouldContinue
	rr.LoopShouldBreak = res.LoopShouldBreak
	return res.Value
}

func (rr *RuntimeResult) Success(value values.BaseValueInterface) *RuntimeResult {
	rr.Reset()
	rr.Value = value
	return rr
}

func (rr *RuntimeResult) SuccessReturn(value values.BaseValueInterface) *RuntimeResult {
	rr.Reset()
	rr.FuncReturnValue = value
	return rr
}

func (rr *RuntimeResult) SuccessContinue() *RuntimeResult {
	rr.Reset()
	rr.LoopShouldContinue = true
	return rr
}

func (rr *RuntimeResult) SuccessBreak() *RuntimeResult {
	rr.Reset()
	rr.LoopShouldBreak = true
	return rr
}

func (rr *RuntimeResult) Failure(err *errors.Error) *RuntimeResult {
	rr.Reset()
	rr.Err = err
	return rr
}

func (rr *RuntimeResult) ShouldReturn() bool {
	return rr.Err != nil || rr.FuncReturnValue != nil || rr.LoopShouldContinue || rr.LoopShouldBreak
}

// *Main Interpreter
func Visit(node nodes.Node, ctx environment.Context) *RuntimeResult {
	switch n := node.(type) {
	case nodes.NumberNode:
		return VisitNumberNode(node.(nodes.NumberNode), ctx)
	case nodes.StringNode:
		return VisitStringNode(node.(nodes.StringNode), ctx)
	case nodes.ListNode:
		return VisitListNode(node.(nodes.ListNode), ctx)
	case nodes.VariableAccessNode:
		return VisitVariableAccessNode(node.(nodes.VariableAccessNode), ctx)
	case nodes.VariableAssignNode:
		return VisitVariableAssignNode(node.(nodes.VariableAssignNode), ctx)
	case nodes.VariableDeleteNode:
		return VisitVariableDeleteNode(node.(nodes.VariableDeleteNode), ctx)
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
	case nodes.ItemAccessNode:
		return VisitItemAccessNode(node.(nodes.ItemAccessNode), ctx)
	case nodes.ItemDeleteNode:
		return VisitItemDeleteNode(node.(nodes.ItemDeleteNode), ctx)
	case nodes.ReturnNode:
		return VisitReturnNode(node.(nodes.ReturnNode), ctx)
	case nodes.ContinueNode:
		return VisitContinueNode(node.(nodes.ContinueNode), ctx)
	case nodes.BreakNode:
		return VisitBreakNode(node.(nodes.BreakNode), ctx)
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
		return res.Failure(errors.NewNotImplementedError(node.Tok.PosRange.Start, node.Tok.PosRange.End, fmt.Sprintf("NumberNode not implemented for token type %s", tokens.TokenTypeNameMap[node.Tok.Type]), ctx))
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

func VisitListNode(node nodes.ListNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	elements := make([]values.BaseValueInterface, 0, len(node.ElementNodes))

	for i := 0; i < len(node.ElementNodes); i++ {
		elements = append(elements, res.Register(Visit(node.ElementNodes[i], ctx)))
		if res.ShouldReturn() {
			return res
		}
	}

	list := &values.List{Elements: elements}
	list.SetContext(ctx)
	list.SetValuePos(node.GetPosRange())
	return res.Success(list)
}

func VisitVariableAccessNode(node nodes.VariableAccessNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	posRange := node.GetPosRange()

	varName := node.VarNameTok.Value.(string)
	rawValue := ctx.SymTable.GetSymbol(varName)

	if rawValue == nil {
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is not defined", varName), ctx))
	}

	// value := rawValue.(values.BaseValueInterface).Copy() //BUG: breaks list operations
	value := rawValue.(values.BaseValueInterface)

	value.SetContext(ctx)
	value.SetValuePos(posRange)
	return res.Success(value)
}

func VisitVariableAssignNode(node nodes.VariableAssignNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	varName := node.VarNameTok.Value.(string)
	value := res.Register(Visit(node.ValueNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	ctx.SymTable.SetSymbol(varName, value)
	return res.Success(value)
}

func VisitVariableDeleteNode(node nodes.VariableDeleteNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	posRange := node.GetPosRange()

	varName := node.VarNameTok.Value.(string)
	rawValue := ctx.SymTable.GetSymbol(varName)

	if rawValue == nil {
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is not defined", varName), ctx))
	}

	ctx.SymTable.RemoveSymbol(varName)

	value := rawValue.(values.BaseValueInterface)

	value.SetContext(ctx)
	value.SetValuePos(posRange)
	return res.Success(value)
}

func VisitBinOpNode(node nodes.BinOpNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	left := res.Register(Visit(node.LeftNode, ctx))
	if res.ShouldReturn() {
		return res
	}
	right := res.Register(Visit(node.RightNode, ctx))
	if res.ShouldReturn() {
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
	if res.ShouldReturn() {
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
		if res.ShouldReturn() {
			return res
		}

		if conditionValue.IsTrue() {
			exprValue := res.Register(Visit(node.Cases[i].Expr, ctx))
			if res.ShouldReturn() {
				return res
			}
			if node.Cases[i].ShouldReturnNull == true {
				goto returnNull
			}
			return res.Success(exprValue)
		}
	}
	if node.ElseCase != nil {
		elseValue := res.Register(Visit(node.ElseCase.Expr, ctx))
		if res.ShouldReturn() {
			return res
		}
		if node.ElseCase.ShouldReturnNull == true {
			goto returnNull
		}
		return res.Success(elseValue)
	}

returnNull:
	return res.Success(&values.Null{})
}

func VisitForNode(node nodes.ForNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	var elements []values.BaseValueInterface

	startValue := res.Register(Visit(node.StartValueNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	stopValue := res.Register(Visit(node.StopValueNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	var stepValue values.BaseValueInterface = &values.Integer{Value: 1}
	if node.StepValueNode != nil {
		stepValue = res.Register(Visit(node.StepValueNode, ctx))
		if res.ShouldReturn() {
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

		element := res.Register(Visit(node.BodyNode, ctx))
		if res.ShouldReturn() && res.LoopShouldContinue == false && res.LoopShouldBreak == false {
			return res
		}
		if res.LoopShouldContinue {
			continue
		}
		if res.LoopShouldBreak {
			break
		}

		if node.ShouldReturnNull == false {
			elements = append(elements, element)
		} //prevents unnecessary alloc

		condRes, err = cond()
		if err != nil {
			return res.Failure(err)
		}
	}

	if node.ShouldReturnNull == true {
		return res.Success(&values.Null{})
	}
	result := &values.List{Elements: elements}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitWhileNode(node nodes.WhileNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	var elements []values.BaseValueInterface

	for {
		condition := res.Register(Visit(node.CondNode, ctx))
		if res.ShouldReturn() {
			return res
		}
		if !condition.IsTrue() {
			break
		}
		element := res.Register(Visit(node.BodyNode, ctx))
		if res.ShouldReturn() && res.LoopShouldContinue == false && res.LoopShouldBreak == false {
			return res
		}
		if res.LoopShouldContinue {
			continue
		}
		if res.LoopShouldBreak {
			break
		}

		if node.ShouldReturnNull == false {
			elements = append(elements, element)
		} //prevents unnecessary alloc
	}

	if node.ShouldReturnNull == true {
		return res.Success(&values.Null{})
	}
	result := &values.List{Elements: elements}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitFuncDefNode(node nodes.FuncDefNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	var funcName *string = nil
	if node.VarNameTok != nil {
		funcNameStrVal := node.VarNameTok.Value.(string)
		funcName = &funcNameStrVal
	}

	bodyNode := node.BodyNode
	argNames := make([]string, 0, len(node.ArgNameToks))
	for i := 0; i < len(node.ArgNameToks); i++ {
		argNames = append(argNames, node.ArgNameToks[i].Value.(string))
	}
	// funcValue := values.NewFunction(funcName, bodyNode, argNames)
	funcValue := &values.Function{BodyNode: bodyNode, ArgNames: argNames, ShouldAutoReturn: node.ShouldAutoReturn, BaseFunction: values.BaseFunction{Name: funcName}}
	funcValue.SetContext(ctx)
	funcValue.SetValuePos(node.GetPosRange())

	if funcName != nil {
		ctx.SymTable.SetSymbol(*funcName, funcValue)
	}
	return res.Success(funcValue)
}

func VisitCallNode(node nodes.CallNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	valueToCall := res.Register(Visit(node.NodeToCall, ctx))
	if res.ShouldReturn() {
		return res
	}
	valueToCall = valueToCall.Copy()
	valueToCall.SetValuePos(node.GetPosRange())

	args := make([]values.BaseValueInterface, 0, len(node.ArgNodes))
	for i := 0; i < len(node.ArgNodes); i++ {
		args = append(args, res.Register(Visit(node.ArgNodes[i], ctx)))
		if res.ShouldReturn() {
			return res
		}
	}
	returnValue := res.Register(ExecuteCallable(valueToCall, args, node, ctx))
	if res.ShouldReturn() {
		return res
	}
	returnValue.SetContext(ctx)
	returnValue.SetValuePos(node.GetPosRange())
	return res.Success(returnValue)
}

func VisitItemAccessNode(node nodes.ItemAccessNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	valueToAccess := res.Register(Visit(node.NodeToAccess, ctx))
	if res.ShouldReturn() {
		return res
	}
	key := res.Register(Visit(node.KeyNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	result, error := valueToAccess.GetItem(key)

	if error != nil {
		return res.Failure(error)
	}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitItemDeleteNode(node nodes.ItemDeleteNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	valueToAccess := res.Register(Visit(node.NodeToAccess, ctx))
	if res.ShouldReturn() {
		return res
	}
	key := res.Register(Visit(node.KeyNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	result, error := valueToAccess.DelItem(key)

	if error != nil {
		return res.Failure(error)
	}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitReturnNode(node nodes.ReturnNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	var value values.BaseValueInterface = &values.Null{}
	if node.NodeToReturn != nil {
		value = res.Register(Visit(node.NodeToReturn, ctx))
		if res.ShouldReturn() {
			return res
		}
	}
	return res.SuccessReturn(value)
}

func VisitContinueNode(node nodes.ContinueNode, ctx environment.Context) *RuntimeResult {
	return NewRuntimeResult().SuccessContinue()
}

func VisitBreakNode(node nodes.BreakNode, ctx environment.Context) *RuntimeResult {
	return NewRuntimeResult().SuccessBreak()
}

// *Built-in Functions
type BuiltInFunctionData struct {
	FunctionRef func(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult
	Args        []string
}

var BuiltInFunctionTable = map[string]BuiltInFunctionData{
	"print": {
		FunctionRef: ExecutePrint,
		Args:        []string{"value"},
	},
	"input": {
		FunctionRef: ExecuteInput,
		Args:        []string{},
	},
	"clear": {
		FunctionRef: ExecuteClear,
		Args:        []string{},
	},
	"typeof": {
		FunctionRef: ExecuteTypeOf,
		Args:        []string{"value"},
	},
	"len": {
		FunctionRef: ExecuteLen,
		Args:        []string{"list"},
	},
	"append": {
		FunctionRef: ExecuteAppend,
		Args:        []string{"list", "value"},
	},
	"pop": {
		FunctionRef: ExecutePop,
		Args:        []string{"list", "idx"},
	},
	"extend": {
		FunctionRef: ExecuteExtend,
		Args:        []string{"list_a", "list_b"},
	},
	"exit": {
		FunctionRef: ExecuteExit,
		Args:        []string{"code"},
	},
}

func ExecutePrint(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	fmt.Println(execCtx.SymTable.GetSymbol("value"))
	return NewRuntimeResult().Success(&values.Null{})
}

func ExecuteInput(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	var input string
	fmt.Scanln(&input)
	return NewRuntimeResult().Success(&values.String{Value: input})
}

func ExecuteClear(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	screen.Clear()
	screen.MoveTopLeft()
	return NewRuntimeResult().Success(&values.Null{})
}

func ExecuteTypeOf(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	value := execCtx.SymTable.GetSymbol("value")
	return NewRuntimeResult().Success(&values.String{Value: reflect.TypeOf(value).Elem().Name()})
}

func ExecuteLen(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	value := execCtx.SymTable.GetSymbol("list").(values.BaseValueInterface)

	length, err := value.Length()
	if err != nil {
		return res.Failure(err)
	}

	return res.Success(length)
}

func ExecuteAppend(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	list := execCtx.SymTable.GetSymbol("list").(values.BaseValueInterface)
	value := execCtx.SymTable.GetSymbol("value").(values.BaseValueInterface)

	switch l := list.(type) {
	case *values.List:
		l.Elements = append(l.Elements, value)
	default:
		// posRange := callable.GetPosRange()
		posRange := list.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be List", execCtx))
	}

	return res.Success(&values.Null{})
}

func ExecutePop(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	// posRange := callable.GetPosRange()

	listValue := execCtx.SymTable.GetSymbol("list").(values.BaseValueInterface)
	list, ok := listValue.(*values.List)
	if ok == false {
		posRange := listValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be List", execCtx))
	}
	idxValue := execCtx.SymTable.GetSymbol("idx").(values.BaseValueInterface)
	idx, ok := idxValue.(*values.Integer)
	if ok == false {
		posRange := idxValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be Integer", execCtx))
	}

	element, err := list.DelItem(idx)
	if err != nil {
		return res.Failure(err)
	}

	return res.Success(element)
}

func ExecuteExtend(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	// posRange := callable.GetPosRange()

	listAValue := execCtx.SymTable.GetSymbol("list_a").(values.BaseValueInterface)
	listA, ok := listAValue.(*values.List)
	if ok == false {
		posRange := listAValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be List", execCtx))
	}
	listBValue := execCtx.SymTable.GetSymbol("list_b").(values.BaseValueInterface)
	listB, ok := listBValue.(*values.List)
	if ok == false {
		posRange := listBValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be List", execCtx))
	}

	listA.Elements = append(listA.Elements, listB.Elements...)

	return res.Success(&values.Null{})
}

func ExecuteExit(callable values.BaseFunctionInterface, execCtx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	codeValue := execCtx.SymTable.GetSymbol("code").(values.BaseValueInterface)
	code, ok := codeValue.(*values.Integer)
	if ok == false {
		posRange := codeValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be Integer", execCtx))
	}

	os.Exit(int(code.Value))

	return res.Success(&values.Null{}) // won't matter anyways
}

// *Function Calls
func CheckArgs(callable values.BaseFunctionInterface, argNames []string, args []values.BaseValueInterface) *RuntimeResult {
	res := NewRuntimeResult()
	posRange := callable.GetPosRange()

	if len(args) > len(argNames) {
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("%d too many args passed into '%s'", len(args)-len(argNames), callable.DisplayName()), callable.GetContext()))
	}
	if len(args) < len(argNames) {
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("%d too few args passed into '%s'", len(argNames)-len(args), callable.DisplayName()), callable.GetContext()))
	}

	return res.Success(nil)
}

func PopulateArgs(callable values.BaseFunctionInterface, argNames []string, args []values.BaseValueInterface, execCtx environment.Context) {
	for i := 0; i < len(args); i++ {
		argName := argNames[i]
		argValue := args[i]
		argValue.SetContext(execCtx)
		execCtx.SymTable.SetSymbol(argName, argValue)
	}
}

func CheckAndPopulateArgs(callable values.BaseFunctionInterface, argNames []string, args []values.BaseValueInterface, execCtx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	res.Register(CheckArgs(callable, argNames, args))
	if res.ShouldReturn() {
		return res
	}
	PopulateArgs(callable, argNames, args, execCtx)
	return res.Success(nil)
}

func ExecuteCallable(callableValue values.BaseValueInterface, args []values.BaseValueInterface, callNode nodes.CallNode, ctx environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	callable, ok := callableValue.(values.BaseFunctionInterface)
	if ok == false {
		goto invalidType
	}

	switch c := callable.(type) {
	case *values.Function:
		execCtx := c.GenerateNewContext()

		res.Register(CheckAndPopulateArgs(c, c.ArgNames, args, execCtx))
		if res.ShouldReturn() {
			return res
		}

		value := res.Register(Visit(c.BodyNode, execCtx))
		if res.ShouldReturn() && res.FuncReturnValue == nil {
			return res
		}

		var retValue values.BaseValueInterface = &values.Null{}
		if c.ShouldAutoReturn == true {
			retValue = value
		}
		if res.FuncReturnValue != nil {
			retValue = res.FuncReturnValue
		}
		return res.Success(retValue)
	case *values.BuiltInFunction:
		execCtx := c.GenerateNewContext()

		functionData, ok := BuiltInFunctionTable[c.DisplayName()]
		if ok == false {
			log.Fatalf("No built-in function '%s' defined", c.Name)
		}

		res.Register(CheckAndPopulateArgs(c, functionData.Args, args, execCtx))
		if res.ShouldReturn() {
			return res
		}

		value := res.Register(functionData.FunctionRef(c, execCtx))
		if res.ShouldReturn() {
			return res
		}
		return res.Success(value)
	default:
		goto invalidType
	}
invalidType:
	nodePosRange := callNode.GetPosRange()
	return res.Failure(errors.NewRuntimeError(nodePosRange.Start, nodePosRange.End, "Illegal operation", ctx))
}
