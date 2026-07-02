package values

import (
	"strconv"

	"github.com/DGTV11/weh-script/context"
	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/position"
)

// *BaseValue and interface
type BaseValueInterface interface {
	GetPosRange() position.PositionRange
	SetValuePos(position.PositionRange)
	GetContext() context.Context
	SetContext(ctx context.Context)
	Add(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Div(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	String() string
}

type BaseValue struct {
	PosRange position.PositionRange
	Ctx      context.Context
}

func (bv *BaseValue) GetPosRange() position.PositionRange {
	return bv.PosRange
}
func (bv *BaseValue) SetValuePos(posRange position.PositionRange) {
	bv.PosRange = posRange
}

func (bv *BaseValue) GetContext() context.Context {
	return bv.Ctx
}
func (bv *BaseValue) SetContext(ctx context.Context) {
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

func (self *Float) String() string {
	return strconv.FormatFloat(self.Value, 'g', -1, 64)
}
