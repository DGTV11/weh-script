package values

import (
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/stanNthe5/stringbuf"

	"github.com/DGTV11/weh-script/environment"
	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/position"
)

// *Helper functions
func IPow(a, n int64) int64 {
	var ans int64 = 1

	for n > 0 {
		last_bit := (n & 1)

		if last_bit != 0 {
			ans = ans * a
		}
		a = a * a

		n = n >> 1
	}
	return ans

} //https://www.geeksforgeeks.org/dsa/fast-exponention-using-bit-manipulation/

func Bool2int64(b bool) int64 {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int64
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
} //https://dev.to/chigbeef_77/bool-int-but-stupid-in-go-3jb3

// *BaseValue and interface
type BaseValueInterface interface {
	GetPosRange() position.PositionRange
	SetValuePos(position.PositionRange)
	GetContext() environment.Context
	SetContext(ctx environment.Context)
	Add(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Div(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Pow(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Lt(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Gt(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Lte(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Gte(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	LNot() (BaseValueInterface, *errors.Error)
	Length() (BaseValueInterface, *errors.Error)
	GetItem(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	DelItem(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	GetAttr(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	DelAttr(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Copy() BaseValueInterface
	IsTrue() bool
	IllegalOperation(other BaseValueInterface) *errors.Error
	String() string
}

type BaseValue struct {
	PosRange position.PositionRange
	Ctx      environment.Context
}

func (bv *BaseValue) GetPosRange() position.PositionRange {
	return bv.PosRange
}
func (bv *BaseValue) SetValuePos(posRange position.PositionRange) {
	bv.PosRange = posRange
}

func (bv *BaseValue) GetContext() environment.Context {
	return bv.Ctx
}
func (bv *BaseValue) SetContext(ctx environment.Context) {
	bv.Ctx = ctx
}

func (self *BaseValue) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Div(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Pow(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Lt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Gt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Lte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Gte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}

//	func (self *BaseValue) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
//		return nil, self.IllegalOperation(other)
//	}
//
//	func (self *BaseValue) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
//		return nil, self.IllegalOperation(other)
//	}
//
//	func (self *BaseValue) LNot() (BaseValueInterface, *errors.Error) {
//		return nil, self.IllegalOperation(nil)
//	}
func (self *BaseValue) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *BaseValue) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *BaseValue) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *BaseValue) Length() (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(nil)
}
func (self *BaseValue) GetItem(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) DelItem(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) GetAttr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) DelAttr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) Copy() BaseValueInterface {
	log.Fatal("No Copy method defined")
	return nil
}

func (self *BaseValue) IsTrue() bool {
	return false
}

func (self *BaseValue) IllegalOperation(other BaseValueInterface) *errors.Error {
	var otherPosRange position.PositionRange
	if other == nil {
		otherPosRange = self.GetPosRange()
	} else {
		otherPosRange = other.GetPosRange()
	}

	return errors.NewRuntimeError(self.GetPosRange().Start, otherPosRange.End, "Illegal operation", self.GetContext())
}

// *Null
type Null struct {
	BaseValue
}

func (self *Null) Copy() BaseValueInterface {
	return self
}
func (self *Null) IsTrue() bool {
	return false
}
func (self *Null) String() string {
	return "null"
}

// *Integer
type Integer struct {
	BaseValue
	Value int64
}

func (self *Integer) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: self.Value + o.Value}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Float{Value: float64(self.Value) + o.Value}
		res.SetContext(self.GetContext())
	default:
		return nil, self.IllegalOperation(other)
	}
	return res, nil
}
func (self *Integer) Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: self.Value - o.Value}
	case *Float:
		res = &Float{Value: float64(self.Value) - o.Value}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: self.Value * o.Value}
	case *Float:
		res = &Float{Value: float64(self.Value) * o.Value}
	case *String:
		if ((len(o.Value) * int(self.Value)) / int(self.Value)) != len(o.Value) { //*detects integer overflow if any, based on https://www.geeksforgeeks.org/dsa/check-integer-overflow-multiplication/
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, other.GetPosRange().End, "String length limit exceeded", self.GetContext())
		}
		res = &String{Value: strings.Repeat(o.Value, int(self.Value))}
	case *List:
		if ((len(o.Elements) * int(self.Value)) / int(self.Value)) != len(o.Elements) { //*detects integer overflow if any, based on https://www.geeksforgeeks.org/dsa/check-integer-overflow-multiplication/
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, other.GetPosRange().End, "Integer length limit exceeded", self.GetContext())
		}
		res = &List{Elements: slices.Repeat(o.Elements, int(self.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Div(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil
	otherPosRange := other.GetPosRange()

	switch o := other.(type) {
	case *Integer:
		if o.Value == 0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero", self.GetContext())
		}
		res = &Float{Value: float64(self.Value) / float64(o.Value)}
	case *Float:
		if o.Value == 0.0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero", self.GetContext())
		}
		res = &Float{Value: float64(self.Value) / o.Value}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Pow(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: IPow(self.Value, o.Value)}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Float{Value: math.Pow(float64(self.Value), o.Value)}
		res.SetContext(self.GetContext())
	default:
		return nil, self.IllegalOperation(other)
	}
	return res, nil
}
func (self *Integer) Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value == o.Value)}
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) == o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value != o.Value)}
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) != o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Lt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value < o.Value)}
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) < o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Gt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value > o.Value)}
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) > o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Lte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value <= o.Value)}
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) <= o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) Gte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value >= o.Value)}
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) >= o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}

//	func (self *Integer) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
//		res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
//		return res, nil
//	}
//
//	func (self *Integer) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
//		res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
//		return res, nil
//	}
//
//	func (self *Integer) LNot() (BaseValueInterface, *errors.Error) {
//		res := &Integer{Value: Bool2int64(!self.IsTrue())}
//		return res, nil
//	}
func (self *Integer) Copy() BaseValueInterface {
	copy := &Integer{Value: self.Value}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *Integer) IsTrue() bool {
	return self.Value != 0
}
func (self *Integer) String() string {
	return strconv.FormatInt(self.Value, 10)
}

// *Float
type Float struct {
	BaseValue
	Value float64
}

func (self *Float) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Float{Value: self.Value + o.Value}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Float{Value: self.Value + float64(o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Float{Value: self.Value - o.Value}
	case *Integer:
		res = &Float{Value: self.Value - float64(o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Float{Value: self.Value * o.Value}
	case *Integer:
		res = &Float{Value: self.Value * float64(o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Div(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil
	otherPosRange := other.GetPosRange()

	switch o := other.(type) {
	case *Float:
		if o.Value == 0.0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero", self.GetContext())
		}
		res = &Float{Value: self.Value / o.Value}
	case *Integer:
		if o.Value == 0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero", self.GetContext())
		}
		res = &Float{Value: self.Value / float64(o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Pow(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Float{Value: math.Pow(self.Value, o.Value)}
	case *Integer:
		res = &Float{Value: math.Pow(self.Value, float64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value == o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value == float64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value != o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value != float64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Lt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value < o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value < float64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Gt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value > o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value > float64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Lte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value <= o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value <= float64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) Gte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value >= o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value >= float64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}

//	func (self *Float) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
//		res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
//		return res, nil
//	}
//
//	func (self *Float) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
//		res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
//		return res, nil
//	}
//
//	func (self *Float) LNot() (BaseValueInterface, *errors.Error) {
//		res := &Integer{Value: Bool2int64(!self.IsTrue())}
//		return res, nil
//	}
func (self *Float) Copy() BaseValueInterface {
	copy := &Float{Value: self.Value}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *Float) IsTrue() bool {
	return self.Value != 0.0
}
func (self *Float) String() string {
	return strconv.FormatFloat(self.Value, 'g', -1, 64)
}

// *String
type String struct {
	BaseValue
	Value string
}

func (self *String) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *String:
		res = &String{Value: self.Value + o.Value}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *String) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		if ((len(self.Value) * int(o.Value)) / int(o.Value)) != len(self.Value) { //*detects integer overflow if any, based on https://www.geeksforgeeks.org/dsa/check-integer-overflow-multiplication/
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, other.GetPosRange().End, "String length limit exceeded", self.GetContext())
		}
		res = &String{Value: strings.Repeat(self.Value, int(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *String) Copy() BaseValueInterface {
	copy := &String{Value: strings.Clone(self.Value)}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *String) IsTrue() bool {
	return len(self.Value) > 0
}
func (self *String) String() string {
	return strconv.Quote(self.Value)
}

// *List
type List struct {
	BaseValue
	Elements []BaseValueInterface
}

func (self *List) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *List:
		res = self.Copy()
		res.(*List).Elements = append(res.(*List).Elements, o.Elements...)
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *List) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		if ((len(self.Elements) * int(o.Value)) / int(o.Value)) != len(self.Elements) { //*detects integer overflow if any, based on https://www.geeksforgeeks.org/dsa/check-integer-overflow-multiplication/
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, other.GetPosRange().End, "Integer length limit exceeded", self.GetContext())
		}
		res = &List{Elements: slices.Repeat(self.Elements, int(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *List) Length() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: int64(len(self.Elements))}
	return res, nil
}
func (self *List) GetItem(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		var idx int
		if rawIdx < 0 {
			idx = len(self.Elements) + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= len(self.Elements) || idx < 0 {
			endPos := other.GetPosRange().End
			x := ' '
			endPos.Advance(&x) //*evil hack
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
		}
		res = self.Elements[idx]
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *List) DelItem(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		var idx int
		if rawIdx < 0 {
			idx = len(self.Elements) + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= len(self.Elements) || idx < 0 {
			endPos := other.GetPosRange().End
			x := ' '
			endPos.Advance(&x) //*evil hack
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be removed from List because index is out of bounds", rawIdx), self.GetContext())
		}
		res = self.Elements[idx]
		self.Elements = append(self.Elements[:idx], self.Elements[idx+1:]...)
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *List) Copy() BaseValueInterface {
	copiedElements := make([]BaseValueInterface, len(self.Elements))
	copy(copiedElements, self.Elements)

	copy := &List{Elements: copiedElements}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *List) IsTrue() bool {
	return len(self.Elements) > 0
}
func (self *List) String() string {
	// return fmt.Sprintf("%v", self.Elements)
	sb := stringbuf.New("[")
	for i := 0; i < len(self.Elements)-1; i++ {
		sb.Append(self.Elements[i].String(), ", ")
	}
	if len(self.Elements) > 0 {
		sb.Append(self.Elements[len(self.Elements)-1].String())
	}
	sb.Append("]")
	return sb.String()
}

// *Function
type Function struct {
	BaseValue
	Name     string
	BodyNode nodes.Node
	ArgNames []string
}

func NewFunction(name *string, bodyNode nodes.Node, argNames []string) *Function {
	var funcName string
	if name == nil {
		funcName = "<anonymous>"
	} else {
		funcName = *name
	}
	return &Function{
		Name:     funcName,
		BodyNode: bodyNode,
		ArgNames: argNames,
	}
}

func (self *Function) String() string {
	return fmt.Sprintf("<function %s>", self.Name)
}

func (self *Function) Copy() BaseValueInterface {
	copy := &Function{Name: self.Name, BodyNode: self.BodyNode, ArgNames: self.ArgNames}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
