package values

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

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
	GetContext() *environment.Context
	SetContext(ctx *environment.Context)
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
	GetItem(key BaseValueInterface) (BaseValueInterface, *errors.Error)
	SetItem(key BaseValueInterface, value BaseValueInterface) (BaseValueInterface, *errors.Error)
	DelItem(key BaseValueInterface) (BaseValueInterface, *errors.Error)
	GetMember(fieldName string, posRange position.PositionRange) (BaseValueInterface, *errors.Error)
	SetMember(fieldName string, value BaseValueInterface, posRange position.PositionRange) (BaseValueInterface, *errors.Error)
	DelMember(fieldName string, posRange position.PositionRange) (BaseValueInterface, *errors.Error)
	Copy() BaseValueInterface
	IsTrue() bool
	IllegalOperation(other BaseValueInterface) *errors.Error
	String() string
	GoString() string
}

type BaseValue struct {
	PosRange position.PositionRange
	Ctx      *environment.Context
}

func (bv *BaseValue) GetPosRange() position.PositionRange {
	return bv.PosRange
}
func (bv *BaseValue) SetValuePos(posRange position.PositionRange) {
	bv.PosRange = posRange
}

func (bv *BaseValue) GetContext() *environment.Context {
	return bv.Ctx
}
func (bv *BaseValue) SetContext(ctx *environment.Context) {
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
func (self *BaseValue) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}

func (self *BaseValue) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(other)
}
func (self *BaseValue) LNot() (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(nil)
}
func (self *BaseValue) Length() (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(nil)
}
func (self *BaseValue) GetItem(key BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(key)
}
func (self *BaseValue) SetItem(key BaseValueInterface, value BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(key)
}
func (self *BaseValue) DelItem(key BaseValueInterface) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperation(key)
}
func (self *BaseValue) GetMember(fieldName string, posRange position.PositionRange) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperationManualPosRange(posRange.Start, posRange.End)
}
func (self *BaseValue) SetMember(fieldName string, value BaseValueInterface, posRange position.PositionRange) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperationManualPosRange(posRange.Start, posRange.End)
}
func (self *BaseValue) DelMember(fieldName string, posRange position.PositionRange) (BaseValueInterface, *errors.Error) {
	return nil, self.IllegalOperationManualPosRange(posRange.Start, posRange.End)
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

	return self.IllegalOperationManualPosRange(self.GetPosRange().Start, otherPosRange.End)
}

func (self *BaseValue) IllegalOperationManualPosRange(posStart *position.Position, posEnd *position.Position) *errors.Error {
	return errors.NewRuntimeError(posStart, posEnd, "Illegal operation", self.GetContext())
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
func (self *Null) GoString() string {
	return self.String()
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
	case *Char:
		res = &Integer{Value: self.Value + int64(o.Value)}
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
	case *Char:
		res = &Integer{Value: self.Value - int64(o.Value)}
		res.SetContext(self.GetContext())
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
	case *Char:
		res = &Integer{Value: self.Value * int64(o.Value)}
		res.SetContext(self.GetContext())
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
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value == int64(o.Value))}
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
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value != int64(o.Value))}
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
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value < int64(o.Value))}
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
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value > int64(o.Value))}
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
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value <= int64(o.Value))}
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
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value >= int64(o.Value))}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *Integer) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *Integer) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
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
func (self *Integer) GoString() string {
	return self.String()
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
	case *Integer:
		res = &Float{Value: self.Value + float64(o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
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
func (self *Float) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *Float) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *Float) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
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
	s := strconv.FormatFloat(self.Value, 'g', -1, 64)
	if !strings.ContainsAny(s, ".eEnN") {
		s += ".0"
	}
	return s
}
func (self *Float) GoString() string {
	return self.String()
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
func (self *String) Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *String:
		res = &Integer{Value: Bool2int64(self.Value == o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *String) Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *String:
		res = &Integer{Value: Bool2int64(self.Value != o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *String) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *String) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *String) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *String) Length() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: int64(utf8.RuneCountInString(self.Value))}
	return res, nil
}
func (self *String) GetItem(key BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := key.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		runeCount := utf8.RuneCountInString(self.Value)
		var idx int
		if rawIdx < 0 {
			idx = runeCount + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= runeCount || idx < 0 {
			// endPos := key.GetPosRange().End
			// x := ' '
			// endPos.Advance(&x) //*evil hack
			// return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, key.GetPosRange().End, fmt.Sprintf("Char at index %d could not be retrieved from String because index is out of bounds", rawIdx), self.GetContext())
		}
		currentIdx := 0
		for _, r := range self.Value {
			if currentIdx == idx {
				res = &Char{Value: r}
				break
			}
			currentIdx++
		} //O(idx) time complexity, should consider using []rune for storage instead
	default:
		return nil, self.IllegalOperation(key)
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
	return self.Value
}
func (self *String) GoString() string {
	return strconv.Quote(self.String())
}

// *Char
type Char struct {
	BaseValue
	Value rune
}

func (self *Char) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Char{Value: self.Value + o.Value}
	case *Integer:
		res = &Char{Value: self.Value + rune(o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Char{Value: self.Value - o.Value}
	case *Integer:
		res = &Char{Value: self.Value - rune(o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Char{Value: self.Value * o.Value}
	case *Integer:
		res = &Char{Value: self.Value * rune(o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value == o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(int64(self.Value) == o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value != o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(int64(self.Value) != o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Lt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value < o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(int64(self.Value) < o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Gt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value > o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(int64(self.Value) > o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Lte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value <= o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(int64(self.Value) <= o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) Gte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Char:
		res = &Integer{Value: Bool2int64(self.Value >= o.Value)}
	case *Integer:
		res = &Integer{Value: Bool2int64(int64(self.Value) >= o.Value)}
	default:
		return nil, self.IllegalOperation(other)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Char) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *Char) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *Char) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *Char) Copy() BaseValueInterface {
	copy := &Char{Value: self.Value}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *Char) IsTrue() bool {
	return self.Value != 0
}
func (self *Char) String() string {
	return string(self.Value)
}
func (self *Char) GoString() string {
	return strconv.QuoteRune(self.Value)
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
func (self *List) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *List) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *List) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *List) Length() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: int64(len(self.Elements))}
	return res, nil
}
func (self *List) GetItem(key BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := key.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		var idx int
		if rawIdx < 0 {
			idx = len(self.Elements) + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= len(self.Elements) || idx < 0 {
			// endPos := key.GetPosRange().End
			// x := ' '
			// endPos.Advance(&x) //*evil hack
			// return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, key.GetPosRange().End, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
		}
		res = self.Elements[idx]
	default:
		return nil, self.IllegalOperation(key)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *List) SetItem(key BaseValueInterface, value BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := key.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		var idx int
		if rawIdx < 0 {
			idx = len(self.Elements) + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= len(self.Elements) || idx < 0 {
			// endPos := key.GetPosRange().End
			// x := ' '
			// endPos.Advance(&x) //*evil hack
			// return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, key.GetPosRange().End, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
		}
		self.Elements[idx] = value
		res = self.Elements[idx]
	default:
		return nil, self.IllegalOperation(key)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *List) DelItem(key BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := key.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		var idx int
		if rawIdx < 0 {
			idx = len(self.Elements) + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= len(self.Elements) || idx < 0 {
			// endPos := key.GetPosRange().End
			// x := ' '
			// endPos.Advance(&x) //*evil hack
			// return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be removed from List because index is out of bounds", rawIdx), self.GetContext())
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, key.GetPosRange().End, fmt.Sprintf("Element at index %d could not be removed from List because index is out of bounds", rawIdx), self.GetContext())
		}
		res = self.Elements[idx]
		self.Elements = append(self.Elements[:idx], self.Elements[idx+1:]...)
	default:
		return nil, self.IllegalOperation(key)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *List) Copy() BaseValueInterface {
	// copiedElements := make([]BaseValueInterface, len(self.Elements))
	// copy(copiedElements, self.Elements)
	//
	// copy := &List{Elements: copiedElements}
	copy := &List{Elements: self.Elements}
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
func (self *List) GoString() string {
	return self.String()
}

// *BaseFunction
type BaseCallableInterface interface {
	BaseValueInterface
	DisplayName() string
}

type BaseFunctionInterface interface {
	BaseCallableInterface
	GenerateNewContext() *environment.Context
}

type BaseFunction struct {
	BaseValue
	Name    *string
	Closure *environment.SymbolTable
}

func (self *BaseFunction) DisplayName() string {
	if self.Name == nil {
		return "<anonymous>"
	}
	return *self.Name
}
func (self *BaseFunction) GenerateNewContext() *environment.Context {
	parentCtx := self.GetContext()
	return &environment.Context{DisplayName: self.DisplayName(), Parent: parentCtx, ParentEntryPos: self.GetPosRange().Start, SymTable: &environment.SymbolTable{Symbols: map[string]any{}, Parent: self.Closure}}
}

// *Function
type Function struct {
	BaseFunction
	BodyNode         nodes.Node
	ArgNames         []string
	ShouldAutoReturn bool
}

func (self *Function) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *Function) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *Function) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *Function) Copy() BaseValueInterface {
	copy := &Function{BodyNode: self.BodyNode, ArgNames: self.ArgNames, ShouldAutoReturn: self.ShouldAutoReturn, BaseFunction: BaseFunction{Name: self.Name, Closure: self.Closure}}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *Function) String() string {
	return fmt.Sprintf("<function %s>", self.DisplayName())
}
func (self *Function) GoString() string {
	return self.String()
}

// *BuiltInFunction
type BuiltInFunction struct {
	BaseFunction
}

func (self *BuiltInFunction) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *BuiltInFunction) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *BuiltInFunction) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *BuiltInFunction) Copy() BaseValueInterface {
	copy := &BuiltInFunction{BaseFunction: BaseFunction{Name: self.Name, Closure: self.Closure}}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *BuiltInFunction) String() string {
	return fmt.Sprintf("<built-in function %s>", self.DisplayName())
}
func (self *BuiltInFunction) GoString() string {
	return self.String()
}

// *File
type File struct {
	BaseValue
	FileValue *os.File
	ModeStr   string
}

func (self *File) Length() (BaseValueInterface, *errors.Error) {
	fInfo, sErr := self.FileValue.Stat()
	if sErr != nil {
		posRange := self.GetPosRange()
		return nil, errors.NewRuntimeError(posRange.Start, posRange.End, sErr.Error(), self.GetContext())
	}
	res := &Integer{Value: fInfo.Size()}
	return res, nil
}
func (self *File) Copy() BaseValueInterface {
	copy := &File{FileValue: self.FileValue, ModeStr: self.ModeStr}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *File) IsTrue() bool {
	return true
}
func (self *File) String() string {
	return fmt.Sprintf("<file path=%s mode=%s>", strconv.Quote(self.FileValue.Name()), strconv.Quote(self.ModeStr))
}
func (self *File) GoString() string {
	return self.String()
}

// *StructDefinition
type StructDefinition struct {
	BaseValue
	Name       *string
	FieldNames []string
}

func (self *StructDefinition) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *StructDefinition) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *StructDefinition) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *StructDefinition) Copy() BaseValueInterface {
	copy := &StructDefinition{Name: self.Name, FieldNames: self.FieldNames}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *StructDefinition) IsTrue() bool {
	return true
}

func (self *StructDefinition) DisplayName() string {
	if self.Name == nil {
		return "<anonymous>"
	}
	return *self.Name
}

func (self *StructDefinition) String() string {
	return fmt.Sprintf("<structure definition %s>", self.DisplayName())
}
func (self *StructDefinition) GoString() string {
	return self.String()
}

// *Structure
type Structure struct {
	BaseValue
	Name            *string
	FieldNameIdxMap map[string]int
	Fields          []BaseValueInterface
}

func (self *Structure) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() && other.IsTrue())}
	return res, nil
}
func (self *Structure) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(self.IsTrue() || other.IsTrue())}
	return res, nil
}
func (self *Structure) LNot() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: Bool2int64(!self.IsTrue())}
	return res, nil
}
func (self *Structure) Length() (BaseValueInterface, *errors.Error) {
	res := &Integer{Value: int64(len(self.Fields))}
	return res, nil
}
func (self *Structure) GetItem(key BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := key.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		var idx int
		if rawIdx < 0 {
			idx = len(self.Fields) + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= len(self.Fields) || idx < 0 {
			// endPos := key.GetPosRange().End
			// x := ' '
			// endPos.Advance(&x) //*evil hack
			// return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, key.GetPosRange().End, fmt.Sprintf("Field at index %d could not be retrieved from Structure because index is out of bounds", rawIdx), self.GetContext())
		}
		res = self.Fields[idx]
	default:
		return nil, self.IllegalOperation(key)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Structure) SetItem(key BaseValueInterface, value BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := key.(type) {
	case *Integer:
		rawIdx := int(o.Value)
		var idx int
		if rawIdx < 0 {
			idx = len(self.Fields) + rawIdx
		} else {
			idx = rawIdx
		}

		if idx >= len(self.Fields) || idx < 0 {
			// endPos := key.GetPosRange().End
			// x := ' '
			// endPos.Advance(&x) //*evil hack
			// return nil, errors.NewRuntimeError(self.GetPosRange().Start, endPos, fmt.Sprintf("Element at index %d could not be retrieved from List because index is out of bounds", rawIdx), self.GetContext())
			return nil, errors.NewRuntimeError(self.GetPosRange().Start, key.GetPosRange().End, fmt.Sprintf("Field at index %d could not be retrieved from Structure because index is out of bounds", rawIdx), self.GetContext())
		}
		self.Fields[idx] = value
		res = self.Fields[idx]
	default:
		return nil, self.IllegalOperation(key)
	}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Structure) GetMember(fieldName string, posRange position.PositionRange) (BaseValueInterface, *errors.Error) {
	// var res BaseValueInterface = nil
	idx, ok := self.FieldNameIdxMap[fieldName]
	if ok == false {
		return nil, errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("Field '%s' does not exist within Structure", fieldName), self.GetContext())
	}
	return self.Fields[idx], nil
}
func (self *Structure) SetMember(fieldName string, value BaseValueInterface, posRange position.PositionRange) (BaseValueInterface, *errors.Error) {
	idx, ok := self.FieldNameIdxMap[fieldName]
	if ok == false {
		return nil, errors.NewRuntimeError(posRange.Start, posRange.End, fmt.Sprintf("Field '%s' does not exist within Structure", fieldName), self.GetContext())
	}
	self.Fields[idx] = value
	return value, nil
}

func (self *Structure) Copy() BaseValueInterface {
	copy := &Structure{Fields: self.Fields, FieldNameIdxMap: self.FieldNameIdxMap}
	copy.SetValuePos(self.GetPosRange())
	copy.SetContext(self.GetContext())
	return copy
}
func (self *Structure) IsTrue() bool {
	return true
}
func (self *Structure) DisplayName() string {
	if self.Name == nil {
		return "<anonymous>"
	}
	return *self.Name
}
func (self *Structure) String() string {
	return fmt.Sprintf("<structure %s>", self.DisplayName())
}
func (self *Structure) GoString() string {
	return self.String() //TODO: full human-readable view into elements
}
