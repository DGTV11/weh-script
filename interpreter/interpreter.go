package interpreter

import (
	"fmt"
	"log"
	"maps"
	"math"
	"os"
	"reflect"

	"github.com/inancgumus/screen"

	"github.com/DGTV11/weh-script/environment"
	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/lexer"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/parser"
	"github.com/DGTV11/weh-script/tokens"
	"github.com/DGTV11/weh-script/values"
)

// *Setup
func SetupGlobalSymbolTable() *environment.SymbolTable {
	GlobalSymbolTable := environment.SymbolTable{Symbols: map[string]any{}}

	//*Load constants
	GlobalSymbolTable.SetSymbol("null", &values.Null{})
	GlobalSymbolTable.SetSymbol("true", &values.Integer{Value: 1})
	GlobalSymbolTable.SetSymbol("false", &values.Integer{Value: 0})
	GlobalSymbolTable.SetSymbol("inf", &values.Float{Value: math.Inf(1)})
	GlobalSymbolTable.SetSymbol("nan", &values.Float{Value: math.NaN()})

	//*Load functions
	for funcName := range maps.Keys(BuiltInFunctionTable) {
		GlobalSymbolTable.SetSymbol(funcName, &values.BuiltInFunction{BaseFunction: values.BaseFunction{Name: &funcName}})
	}

	return &GlobalSymbolTable
}

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
func Visit(node nodes.Node, ctx *environment.Context) *RuntimeResult {
	switch n := node.(type) {
	case nodes.NumberNode:
		return VisitNumberNode(node.(nodes.NumberNode), ctx)
	case nodes.StringNode:
		return VisitStringNode(node.(nodes.StringNode), ctx)
	case nodes.CharNode:
		return VisitCharNode(node.(nodes.CharNode), ctx)
	case nodes.ListNode:
		return VisitListNode(node.(nodes.ListNode), ctx)
	case nodes.StatementsNode:
		return VisitStatementsNode(node.(nodes.StatementsNode), ctx)
	case nodes.VariableAccessNode:
		return VisitVariableAccessNode(node.(nodes.VariableAccessNode), ctx)
	case nodes.VariableAssignNode:
		return VisitVariableAssignNode(node.(nodes.VariableAssignNode), ctx)
	case nodes.VariableReassignNode:
		return VisitVariableReassignNode(node.(nodes.VariableReassignNode), ctx)
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
	case nodes.StructDefNode:
		return VisitStructDefNode(node.(nodes.StructDefNode), ctx)
	case nodes.CallNode:
		return VisitCallNode(node.(nodes.CallNode), ctx)
	case nodes.ItemAccessNode:
		return VisitItemAccessNode(node.(nodes.ItemAccessNode), ctx)
	case nodes.ItemAssignNode:
		return VisitItemAssignNode(node.(nodes.ItemAssignNode), ctx)
	case nodes.ItemDeleteNode:
		return VisitItemDeleteNode(node.(nodes.ItemDeleteNode), ctx)
	case nodes.MemberAccessNode:
		return VisitMemberAccessNode(node.(nodes.MemberAccessNode), ctx)
	case nodes.MemberAssignNode:
		return VisitMemberAssignNode(node.(nodes.MemberAssignNode), ctx)
	case nodes.MemberDeleteNode:
		return VisitMemberDeleteNode(node.(nodes.MemberDeleteNode), ctx)
	case nodes.ReturnNode:
		return VisitReturnNode(node.(nodes.ReturnNode), ctx)
	case nodes.ContinueNode:
		return VisitContinueNode(node.(nodes.ContinueNode), ctx)
	case nodes.BreakNode:
		return VisitBreakNode(node.(nodes.BreakNode), ctx)
	case nodes.ImportNode:
		return VisitImportNode(node.(nodes.ImportNode), ctx)
	default:
		posRange := node.GetPosRange()
		return NewRuntimeResult().Failure(errors.NewNotImplementedError(posRange.Start, posRange.End, fmt.Sprintf("No Visit function defined for node type %T", n), ctx))
	}
}

func VisitNumberNode(node nodes.NumberNode, ctx *environment.Context) *RuntimeResult {
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

func VisitStringNode(node nodes.StringNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	str := &values.String{Value: node.Tok.Value.(string)}
	str.SetContext(ctx)
	str.SetValuePos(node.GetPosRange())
	return res.Success(str)
}

func VisitCharNode(node nodes.CharNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	str := &values.Char{Value: node.Tok.Value.(rune)}
	str.SetContext(ctx)
	str.SetValuePos(node.GetPosRange())
	return res.Success(str)
}

func VisitListNode(node nodes.ListNode, ctx *environment.Context) *RuntimeResult {
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

func VisitStatementsNode(node nodes.StatementsNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	var lastStatementValue values.BaseValueInterface = &values.Null{}

	for i := 0; i < len(node.StatementNodes); i++ {
		value := res.Register(Visit(node.StatementNodes[i], ctx))
		if res.ShouldReturn() {
			return res
		}
		if reflect.TypeOf(value).String() != "*values.Null" {
			lastStatementValue = value
		}
	}

	return res.Success(lastStatementValue)
}

func VisitVariableAccessNode(node nodes.VariableAccessNode, ctx *environment.Context) *RuntimeResult {
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

func VisitVariableAssignNode(node nodes.VariableAssignNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	varName := node.VarNameTok.Value.(string)

	var value values.BaseValueInterface = &values.Null{}
	if node.ValueNode != nil {
		value = res.Register(Visit(node.ValueNode, ctx))
		if res.ShouldReturn() {
			return res
		}
	}

	stRes := ctx.SymTable.SetSymbol(varName, value)
	if stRes == false {
		posRange := node.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is already defined", varName), ctx))
	}
	return res.Success(value)
}

func VisitVariableReassignNode(node nodes.VariableReassignNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	varName := node.VarNameTok.Value.(string)
	value := res.Register(Visit(node.ValueNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	stRes := ctx.SymTable.UpdateSymbol(varName, value, node.Nonlocal)
	if stRes == false {
		posRange := node.GetPosRange()
		if node.Nonlocal {
			return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is not defined in parent scope", varName), ctx))
		}
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is not defined", varName), ctx))
	}
	return res.Success(value)
}

func VisitVariableDeleteNode(node nodes.VariableDeleteNode, ctx *environment.Context) *RuntimeResult {
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

func VisitBinOpNode(node nodes.BinOpNode, ctx *environment.Context) *RuntimeResult {
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

func VisitUnaryOpNode(node nodes.UnaryOpNode, ctx *environment.Context) *RuntimeResult {
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

func VisitIfNode(node nodes.IfNode, ctx *environment.Context) *RuntimeResult {
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

func VisitForNode(node nodes.ForNode, ctx *environment.Context) *RuntimeResult {
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
		ctx.SymTable.ForceSetSymbol(node.VarNameTok.Value.(string), i)
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

func VisitWhileNode(node nodes.WhileNode, ctx *environment.Context) *RuntimeResult {
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

func VisitFuncDefNode(node nodes.FuncDefNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	var funcName *string = nil
	if node.VarNameTok != nil {
		funcNameStrVal := node.VarNameTok.Value.(string)
		funcName = &funcNameStrVal
	}

	bodyNode := node.BodyNode
	argNames := make([]string, 0, len(node.ArgNameToks))
	keys := make(map[string]bool)
	for i := 0; i < len(node.ArgNameToks); i++ {
		argName := node.ArgNameToks[i].Value.(string)
		if _, duplicatePresent := keys[argName]; duplicatePresent {
			return res.Failure(errors.NewRuntimeError(node.ArgNameToks[i].PosRange.Start, node.ArgNameToks[i].PosRange.End, fmt.Sprintf("Duplicate arg name '%s'", argName), ctx))
		}
		argNames = append(argNames, argName)
		keys[argName] = true
	}
	// funcValue := values.NewFunction(funcName, bodyNode, argNames)
	funcValue := &values.Function{BodyNode: bodyNode, ArgNames: argNames, ShouldAutoReturn: node.ShouldAutoReturn, BaseFunction: values.BaseFunction{Name: funcName, Closure: ctx.SymTable}}
	funcValue.SetContext(ctx)
	funcValue.SetValuePos(node.GetPosRange())

	if funcName != nil {
		stRes := ctx.SymTable.SetSymbol(*funcName, funcValue)
		if stRes == false {
			posRange := node.GetPosRange()
			return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is already defined", *funcName), ctx))
		}
	}
	// fmt.Println(funcValue.Closure)
	return res.Success(funcValue)
}

func VisitStructDefNode(node nodes.StructDefNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	var structName *string = nil
	if node.VarNameTok != nil {
		structNameStrVal := node.VarNameTok.Value.(string)
		structName = &structNameStrVal
	}

	fieldNames := make([]string, 0, len(node.FieldNameToks))
	keys := make(map[string]bool)
	for i := 0; i < len(node.FieldNameToks); i++ {
		fieldName := node.FieldNameToks[i].Value.(string)
		if _, duplicatePresent := keys[fieldName]; duplicatePresent {
			return res.Failure(errors.NewRuntimeError(node.FieldNameToks[i].PosRange.Start, node.FieldNameToks[i].PosRange.End, fmt.Sprintf("Duplicate field name '%s'", fieldName), ctx))
		}
		fieldNames = append(fieldNames, fieldName)
		keys[fieldName] = true
	}
	structValue := &values.StructDefinition{Name: structName, FieldNames: fieldNames}
	structValue.SetContext(ctx)
	structValue.SetValuePos(node.GetPosRange())

	if structName != nil {
		stRes := ctx.SymTable.SetSymbol(*structName, structValue)
		if stRes == false {
			posRange := node.GetPosRange()
			return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("'%s' is already defined", *structName), ctx))
		}
	}
	return res.Success(structValue)
}

func VisitCallNode(node nodes.CallNode, ctx *environment.Context) *RuntimeResult {
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

func VisitItemAccessNode(node nodes.ItemAccessNode, ctx *environment.Context) *RuntimeResult {
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

func VisitItemAssignNode(node nodes.ItemAssignNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	valueToAssignTo := res.Register(Visit(node.NodeToAssignTo, ctx))
	if res.ShouldReturn() {
		return res
	}
	key := res.Register(Visit(node.KeyNode, ctx))
	if res.ShouldReturn() {
		return res
	}
	value := res.Register(Visit(node.ValueNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	result, error := valueToAssignTo.SetItem(key, value)

	if error != nil {
		return res.Failure(error)
	}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitItemDeleteNode(node nodes.ItemDeleteNode, ctx *environment.Context) *RuntimeResult {
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

func VisitMemberAccessNode(node nodes.MemberAccessNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	valueToAccess := res.Register(Visit(node.NodeToAccess, ctx))
	if res.ShouldReturn() {
		return res
	}

	result, error := valueToAccess.GetMember(node.FieldNameTok.Value.(string), node.GetPosRange())

	if error != nil {
		return res.Failure(error)
	}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitMemberAssignNode(node nodes.MemberAssignNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	valueToAssignTo := res.Register(Visit(node.NodeToAssignTo, ctx))
	if res.ShouldReturn() {
		return res
	}
	value := res.Register(Visit(node.ValueNode, ctx))
	if res.ShouldReturn() {
		return res
	}

	result, error := valueToAssignTo.SetMember(node.FieldNameTok.Value.(string), value, node.GetPosRange())

	if error != nil {
		return res.Failure(error)
	}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitMemberDeleteNode(node nodes.MemberDeleteNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	valueToAccess := res.Register(Visit(node.NodeToAccess, ctx))
	if res.ShouldReturn() {
		return res
	}

	result, error := valueToAccess.DelMember(node.FieldNameTok.Value.(string), node.GetPosRange())

	if error != nil {
		return res.Failure(error)
	}
	result.SetContext(ctx)
	result.SetValuePos(node.GetPosRange())
	return res.Success(result)
}

func VisitReturnNode(node nodes.ReturnNode, ctx *environment.Context) *RuntimeResult {
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

func VisitContinueNode(node nodes.ContinueNode, ctx *environment.Context) *RuntimeResult {
	return NewRuntimeResult().SuccessContinue()
}

func VisitBreakNode(node nodes.BreakNode, ctx *environment.Context) *RuntimeResult {
	return NewRuntimeResult().SuccessBreak()
}

func VisitImportNode(node nodes.ImportNode, ctx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	modulePath := node.ModulePathTok.Value.(string)
	programBytestr, Rerr := os.ReadFile(modulePath)
	if Rerr != nil {
		return res.Failure(errors.NewRuntimeError(node.GetPosRange().Start, node.GetPosRange().End, Rerr.Error(), ctx))
	}

	program := string(programBytestr)

	_lexer := lexer.NewLexer(modulePath, program)
	tokens, Lerr := _lexer.Tokenise()
	if Lerr != nil {
		return res.Failure(Lerr)
	}

	_parser := parser.NewParser(tokens)
	ast := _parser.Parse()
	if ast.Err != nil {
		return res.Failure(ast.Err)
	}

	Iresult := Visit(ast.Node, ctx)
	if Iresult.Err != nil {
		return res.Failure(Iresult.Err)
	}

	return res.Success(&values.Null{})
}

// *Built-in Functions
type BuiltInFunctionData struct {
	FunctionRef func(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult
	Args        []string
}

var BuiltInFunctionTable = map[string]BuiltInFunctionData{
	"print": {
		FunctionRef: ExecutePrint,
		Args:        []string{"value"},
	},
	"println": {
		FunctionRef: ExecutePrintln,
		Args:        []string{"value"},
	},
	// "printf": {
	// 	FunctionRef: ExecutePrintf,
	// 	Args:        []string{"format", "args"},
	// },
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
	"repr": {
		FunctionRef: ExecuteRepr,
		Args:        []string{"value"},
	},
	"len": {
		FunctionRef: ExecuteLen,
		Args:        []string{"list"},
	},
	"hex": {
		FunctionRef: ExecuteHex,
		Args:        []string{"int"},
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
	"fopen": {
		FunctionRef: ExecuteFileOpen,
		Args:        []string{"path", "mode"},
	},
	"fclose": {
		FunctionRef: ExecuteFileClose,
		Args:        []string{"file"},
	},
	"fseek": {
		FunctionRef: ExecuteFileSeek,
		Args:        []string{"file", "offset", "whence"},
	},
	"fread": {
		FunctionRef: ExecuteFileRead,
		Args:        []string{"file", "size"},
	},
	"fwrite": {
		FunctionRef: ExecuteFileWrite,
		Args:        []string{"file", "text"},
	},
	// "fprintf": {
	// 	FunctionRef: ExecuteFilePrintf,
	// 	Args:        []string{"file", "format", "args"},
	// },
	"ftruncate": {
		FunctionRef: ExecuteFileTruncate,
		Args:        []string{"file", "size"},
	},
	"fcreate": {
		FunctionRef: ExecuteFileCreate,
		Args:        []string{"path"},
	},
	"fcreate_temp": {
		FunctionRef: ExecuteFileCreateTemp,
		Args:        []string{"dir", "pattern"},
	},
}

func ExecutePrint(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	fmt.Print(execCtx.SymTable.GetSymbol("value"))
	return NewRuntimeResult().Success(&values.Null{})
}

func ExecutePrintln(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	fmt.Println(execCtx.SymTable.GetSymbol("value"))
	return NewRuntimeResult().Success(&values.Null{})
}

// func ExecutePrintf(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
// 	res := NewRuntimeResult()
//
// 	formatValue := execCtx.SymTable.GetSymbol("format").(values.BaseValueInterface)
// 	format, ok := formatValue.(*values.String)
// 	if ok == false {
// 		posRange := formatValue.GetPosRange()
// 		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be String", execCtx))
// 	}
// 	argsValue := execCtx.SymTable.GetSymbol("args").(values.BaseValueInterface)
// 	args, ok := argsValue.(*values.List)
// 	if ok == false {
// 		posRange := argsValue.GetPosRange()
// 		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be List", execCtx))
// 	}
// 	fmt.Printf(format.Value, args.Elements...)
// 	return NewRuntimeResult().Success(&values.Null{})
// } //TODO: make 'native' printf

func ExecuteInput(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	var input string
	fmt.Scanln(&input)
	return NewRuntimeResult().Success(&values.String{Value: input})
}

func ExecuteClear(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	screen.Clear()
	screen.MoveTopLeft()
	return NewRuntimeResult().Success(&values.Null{})
}

func ExecuteTypeOf(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	value := execCtx.SymTable.GetSymbol("value")
	return NewRuntimeResult().Success(&values.String{Value: reflect.TypeOf(value).Elem().Name()})
}

func ExecuteRepr(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	value := execCtx.SymTable.GetSymbol("value")
	return NewRuntimeResult().Success(&values.String{Value: value.(values.BaseValueInterface).GoString()})
}

func ExecuteLen(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	value := execCtx.SymTable.GetSymbol("list").(values.BaseValueInterface)

	length, err := value.Length()
	if err != nil {
		return res.Failure(err)
	}

	return res.Success(length)
}

func ExecuteHex(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	intValue := execCtx.SymTable.GetSymbol("int").(values.BaseValueInterface)
	int, ok := intValue.(*values.Integer)
	if ok == false {
		posRange := intValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be Integer", execCtx))
	}

	return NewRuntimeResult().Success(&values.String{Value: fmt.Sprintf("0x%x", int.Value)})
}

func ExecuteAppend(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
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

func ExecutePop(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
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

func ExecuteExtend(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
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

func ExecuteExit(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
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

func ExecuteFileOpen(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	pathValue := execCtx.SymTable.GetSymbol("path").(values.BaseValueInterface)
	path, ok := pathValue.(*values.String)
	if ok == false {
		posRange := pathValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be String", execCtx))
	}
	modeValue := execCtx.SymTable.GetSymbol("mode").(values.BaseValueInterface)
	mode, ok := modeValue.(*values.String)
	if ok == false {
		posRange := modeValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be String", execCtx))
	}

	var modeInt int
	switch mode.Value {
	case "r":
		modeInt = os.O_RDONLY
	case "w":
		modeInt = os.O_WRONLY | os.O_CREATE
	case "a":
		modeInt = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	case "r+":
		modeInt = os.O_RDWR
	case "w+":
		modeInt = os.O_RDWR | os.O_CREATE
	case "a+":
		modeInt = os.O_RDWR | os.O_APPEND | os.O_CREATE
	default:
		posRange := modeValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Invalid open mode", execCtx))
	}
	f, err := os.OpenFile(path.Value, modeInt, 0644)
	if err != nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	if f == nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Unknown error in opening file", execCtx))
	}

	if mode.Value[0] == 'w' {
		tErr := f.Truncate(0)
		if tErr != nil {
			posRange := callable.GetPosRange()
			return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, tErr.Error(), execCtx))
		}
	}

	return res.Success(&values.File{FileValue: f, ModeStr: mode.Value})
}

func ExecuteFileClose(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	fileValue := execCtx.SymTable.GetSymbol("file").(values.BaseValueInterface)
	file, ok := fileValue.(*values.File)
	if ok == false {
		posRange := fileValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be File", execCtx))
	}
	err := file.FileValue.Close()
	if err != nil {
		posRange := fileValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	return res.Success(&values.Null{})
}

func ExecuteFileSeek(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	fileValue := execCtx.SymTable.GetSymbol("file").(values.BaseValueInterface)
	file, ok := fileValue.(*values.File)
	if ok == false {
		posRange := fileValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be File", execCtx))
	}
	offsetValue := execCtx.SymTable.GetSymbol("offset").(values.BaseValueInterface)
	offset, ok := offsetValue.(*values.Integer)
	if ok == false {
		posRange := offsetValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be Integer", execCtx))
	}
	whenceValue := execCtx.SymTable.GetSymbol("whence").(values.BaseValueInterface)
	whence, ok := whenceValue.(*values.Integer)
	if ok == false {
		posRange := whenceValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Third argument must be Integer", execCtx))
	}

	newOffset, err := file.FileValue.Seek(offset.Value, int(whence.Value))
	if err != nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	return res.Success(&values.Integer{Value: int64(newOffset)})
}

func ExecuteFileRead(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	fileValue := execCtx.SymTable.GetSymbol("file").(values.BaseValueInterface)
	file, ok := fileValue.(*values.File)
	if ok == false {
		posRange := fileValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be File", execCtx))
	}
	sizeValue := execCtx.SymTable.GetSymbol("size").(values.BaseValueInterface)
	size, ok := sizeValue.(*values.Integer)
	if ok == false {
		posRange := sizeValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be Integer", execCtx))
	}

	var sizeInt int64
	if size.Value < 0 {
		lenValue, lErr := file.Length()
		if lErr != nil {
			return res.Failure(lErr)
		}
		sizeInt = lenValue.(*values.Integer).Value
	} else {
		sizeInt = size.Value
	}

	readBuf := make([]byte, sizeInt)

	noBytesRead, err := file.FileValue.Read(readBuf)
	if err != nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	return res.Success(&values.String{Value: string(readBuf[:noBytesRead])})
}

func ExecuteFileWrite(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	fileValue := execCtx.SymTable.GetSymbol("file").(values.BaseValueInterface)
	file, ok := fileValue.(*values.File)
	if ok == false {
		posRange := fileValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be File", execCtx))
	}
	textValue := execCtx.SymTable.GetSymbol("text").(values.BaseValueInterface)
	text, ok := textValue.(*values.String)
	if ok == false {
		posRange := textValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be String", execCtx))
	}

	noBytesWritten, err := file.FileValue.WriteString(text.Value)
	if err != nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	return res.Success(&values.Integer{Value: int64(noBytesWritten)})
}

func ExecuteFileTruncate(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	fileValue := execCtx.SymTable.GetSymbol("file").(values.BaseValueInterface)
	file, ok := fileValue.(*values.File)
	if ok == false {
		posRange := fileValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be File", execCtx))
	}
	sizeValue := execCtx.SymTable.GetSymbol("size").(values.BaseValueInterface)
	size, ok := sizeValue.(*values.Integer)
	if ok == false {
		posRange := sizeValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be Integer", execCtx))
	}

	err := file.FileValue.Truncate(size.Value)
	if err != nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	return res.Success(&values.Integer{Value: size.Value})
}

func ExecuteFileCreate(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	pathValue := execCtx.SymTable.GetSymbol("path").(values.BaseValueInterface)
	path, ok := pathValue.(*values.String)
	if ok == false {
		posRange := pathValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be String", execCtx))
	}

	f, err := os.Create(path.Value)
	if err != nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	if f == nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Unknown error in creating file", execCtx))
	}

	return res.Success(&values.File{FileValue: f, ModeStr: "r+"})
}

func ExecuteFileCreateTemp(callable values.BaseFunctionInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()

	dirValue := execCtx.SymTable.GetSymbol("dir").(values.BaseValueInterface)
	dir, ok := dirValue.(*values.String)
	if ok == false {
		posRange := dirValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "First argument must be String", execCtx))
	}
	patternValue := execCtx.SymTable.GetSymbol("pattern").(values.BaseValueInterface)
	pattern, ok := patternValue.(*values.String)
	if ok == false {
		posRange := patternValue.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Second argument must be String", execCtx))
	}

	f, err := os.CreateTemp(dir.Value, pattern.Value)
	if err != nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, err.Error(), execCtx))
	}
	if f == nil {
		posRange := callable.GetPosRange()
		return res.Failure(errors.NewRuntimeError(posRange.Start, posRange.End, "Unknown error in creating file", execCtx))
	}

	return res.Success(&values.File{FileValue: f, ModeStr: "r+"})
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

func PopulateArgs(callable values.BaseFunctionInterface, argNames []string, args []values.BaseValueInterface, execCtx *environment.Context) {
	for i := 0; i < len(args); i++ {
		argName := argNames[i]
		argValue := args[i]
		argValue.SetContext(execCtx)
		execCtx.SymTable.SetSymbol(argName, argValue)
	}
}

func CheckAndPopulateArgs(callable values.BaseFunctionInterface, argNames []string, args []values.BaseValueInterface, execCtx *environment.Context) *RuntimeResult {
	res := NewRuntimeResult()
	res.Register(CheckArgs(callable, argNames, args))
	if res.ShouldReturn() {
		return res
	}
	PopulateArgs(callable, argNames, args, execCtx)
	return res.Success(nil)
}

func ExecuteCallable(callableValue values.BaseValueInterface, args []values.BaseValueInterface, callNode nodes.CallNode, ctx *environment.Context) *RuntimeResult {
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
