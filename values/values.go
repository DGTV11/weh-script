package values

import (
	"math"
	"strconv"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/position"
	"github.com/DGTV11/weh-script/runtime"
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
	GetContext() runtime.Context
	SetContext(ctx runtime.Context)
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
	String() string
}

type BaseValue struct {
	PosRange position.PositionRange
	Ctx      runtime.Context
}

func (bv *BaseValue) GetPosRange() position.PositionRange {
	return bv.PosRange
}
func (bv *BaseValue) SetValuePos(posRange position.PositionRange) {
	bv.PosRange = posRange
}

func (bv *BaseValue) GetContext() runtime.Context {
	return bv.Ctx
}
func (bv *BaseValue) SetContext(ctx runtime.Context) {
	bv.Ctx = ctx
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
	}
	return res, nil
}
func (self *Integer) Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: self.Value - o.Value}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Float{Value: float64(self.Value) - o.Value}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: self.Value * o.Value}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Float{Value: float64(self.Value) * o.Value}
		res.SetContext(self.GetContext())
	}
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
		res.SetContext(self.GetContext())
	case *Float:
		if o.Value == 0.0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero", self.GetContext())
		}
		res = &Float{Value: float64(self.Value) / o.Value}
		res.SetContext(self.GetContext())
	}
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
	}
	return res, nil
}
func (self *Integer) Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value == o.Value)}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) == o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value != o.Value)}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) != o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) Lt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value < o.Value)}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) < o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) Gt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value > o.Value)}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) > o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) Lte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value <= o.Value)}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) <= o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) Gte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value >= o.Value)}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64(float64(self.Value) >= o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64((self.Value != 0) && (o.Value != 0))}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64((float64(self.Value) != 0.0) && (o.Value != 0.0))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Integer:
		res = &Integer{Value: Bool2int64((self.Value != 0) || (o.Value != 0))}
		res.SetContext(self.GetContext())
	case *Float:
		res = &Integer{Value: Bool2int64((float64(self.Value) != 0.0) || (o.Value != 0.0))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Integer) LNot() (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	res = &Integer{Value: Bool2int64(self.Value == 0)}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Integer) String() string {
	return strconv.FormatInt(self.Value, 10)
}

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
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Float{Value: self.Value - float64(o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Float{Value: self.Value * o.Value}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Float{Value: self.Value * float64(o.Value)}
		res.SetContext(self.GetContext())
	}
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
		res.SetContext(self.GetContext())
	case *Integer:
		if o.Value == 0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero", self.GetContext())
		}
		res = &Float{Value: self.Value / float64(o.Value)}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Pow(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Float{Value: math.Pow(self.Value, o.Value)}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Float{Value: math.Pow(self.Value, float64(o.Value))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Eq(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value == o.Value)}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value == float64(o.Value))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Ne(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value != o.Value)}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value != float64(o.Value))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Lt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value < o.Value)}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value < float64(o.Value))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Gt(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value > o.Value)}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value > float64(o.Value))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Lte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value <= o.Value)}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value <= float64(o.Value))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) Gte(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64(self.Value >= o.Value)}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64(self.Value >= float64(o.Value))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) LAnd(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64((self.Value != 0.0) && (o.Value != 0.0))}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64((self.Value != 0.0) && (float64(o.Value) != 0.0))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) LOr(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	switch o := other.(type) {
	case *Float:
		res = &Integer{Value: Bool2int64((self.Value != 0.0) || (o.Value != 0.0))}
		res.SetContext(self.GetContext())
	case *Integer:
		res = &Integer{Value: Bool2int64((self.Value != 0.0) || (float64(o.Value) != 0.0))}
		res.SetContext(self.GetContext())
	}
	return res, nil
}
func (self *Float) LNot() (BaseValueInterface, *errors.Error) {
	var res BaseValueInterface = nil

	res = &Integer{Value: Bool2int64(self.Value == 0.0)}
	res.SetContext(self.GetContext())
	return res, nil
}
func (self *Float) String() string {
	return strconv.FormatFloat(self.Value, 'g', -1, 64)
}
