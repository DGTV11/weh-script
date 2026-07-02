package values

import (
	"strconv"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/position"
)

// *BaseValue and interface
type BaseValueInterface interface {
	GetPosRange() position.PositionRange
	SetValuePos(position.PositionRange)
	Add(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	Div(other BaseValueInterface) (BaseValueInterface, *errors.Error)
	String() string
}

type BaseValue struct {
	PosRange position.PositionRange
}

func (bv *BaseValue) GetPosRange() position.PositionRange {
	return bv.PosRange
}
func (bv *BaseValue) SetValuePos(posRange position.PositionRange) {
	bv.PosRange = posRange
}

// *Integer
type Integer struct {
	BaseValue
	Value int64
}

func (self *Integer) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	switch o := other.(type) {
	case *Integer:
		return &Integer{Value: self.Value + o.Value}, nil
	case *Float:
		return &Float{Value: float64(self.Value) + o.Value}, nil
	}
	return nil, nil
}
func (self *Integer) Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	switch o := other.(type) {
	case *Integer:
		return &Integer{Value: self.Value - o.Value}, nil
	case *Float:
		return &Float{Value: float64(self.Value) - o.Value}, nil
	}
	return nil, nil
}
func (self *Integer) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	switch o := other.(type) {
	case *Integer:
		return &Integer{Value: self.Value * o.Value}, nil
	case *Float:
		return &Float{Value: float64(self.Value) * o.Value}, nil
	}
	return nil, nil
}
func (self *Integer) Div(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	otherPosRange := other.GetPosRange()
	switch o := other.(type) {
	case *Integer:
		if o.Value == 0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero")
		}
		return &Float{Value: float64(self.Value) / float64(o.Value)}, nil
	case *Float:
		if o.Value == 0.0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero")
		}
		return &Float{Value: float64(self.Value) / o.Value}, nil
	}
	return nil, nil
}
func (self *Integer) String() string {
	return strconv.FormatInt(self.Value, 10)
}

type Float struct {
	BaseValue
	Value float64
}

func (self *Float) Add(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	switch o := other.(type) {
	case *Float:
		return &Float{Value: self.Value + o.Value}, nil
	case *Integer:
		return &Float{Value: self.Value + float64(o.Value)}, nil
	}
	return nil, nil
}
func (self *Float) Sub(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	switch o := other.(type) {
	case *Float:
		return &Float{Value: self.Value - o.Value}, nil
	case *Integer:
		return &Float{Value: self.Value - float64(o.Value)}, nil
	}
	return nil, nil
}
func (self *Float) Mul(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	switch o := other.(type) {
	case *Float:
		return &Float{Value: self.Value * o.Value}, nil
	case *Integer:
		return &Float{Value: self.Value * float64(o.Value)}, nil
	}
	return nil, nil
}
func (self *Float) Div(other BaseValueInterface) (BaseValueInterface, *errors.Error) {
	otherPosRange := other.GetPosRange()
	switch o := other.(type) {
	case *Float:
		if o.Value == 0.0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero")
		}
		return &Float{Value: self.Value / o.Value}, nil
	case *Integer:
		if o.Value == 0 {
			return nil, errors.NewRuntimeError(otherPosRange.Start, otherPosRange.End, "Division by zero")
		}
		return &Float{Value: self.Value / float64(o.Value)}, nil
	}
	return nil, nil
}

func (self *Float) String() string {
	return strconv.FormatFloat(self.Value, 'g', -1, 64)
}
