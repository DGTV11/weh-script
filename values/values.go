package values

import (
	"strconv"

	"github.com/DGTV11/weh-script/position"
)

// *BaseValue and interface
type BaseValueInterface interface {
	GetPosRange() position.PositionRange
	SetValuePos(position.PositionRange)
	Add(other BaseValueInterface) BaseValueInterface
	Sub(other BaseValueInterface) BaseValueInterface
	Mul(other BaseValueInterface) BaseValueInterface
	Div(other BaseValueInterface) BaseValueInterface
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

func (self *Integer) Add(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Integer:
		return &Integer{Value: self.Value + o.Value}
	case *Float:
		return &Float{Value: float64(self.Value) + o.Value}
	}
	return nil
}
func (self *Integer) Sub(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Integer:
		return &Integer{Value: self.Value - o.Value}
	case *Float:
		return &Float{Value: float64(self.Value) - o.Value}
	}
	return nil
}
func (self *Integer) Mul(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Integer:
		return &Integer{Value: self.Value * o.Value}
	case *Float:
		return &Float{Value: float64(self.Value) * o.Value}
	}
	return nil
}
func (self *Integer) Div(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Integer:
		// return &Integer{Value: self.Value / o.Value}
		return &Float{Value: float64(self.Value) / float64(o.Value)}
	case *Float:
		return &Float{Value: float64(self.Value) / o.Value}
	}
	return nil
}
func (self *Integer) String() string {
	return strconv.FormatInt(self.Value, 10)
}

type Float struct {
	BaseValue
	Value float64
}

func (self *Float) Add(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Float:
		return &Float{Value: self.Value + o.Value}
	case *Integer:
		return &Float{Value: self.Value + float64(o.Value)}
	}
	return nil
}
func (self *Float) Sub(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Float:
		return &Float{Value: self.Value - o.Value}
	case *Integer:
		return &Float{Value: self.Value - float64(o.Value)}
	}
	return nil
}
func (self *Float) Mul(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Float:
		return &Float{Value: self.Value * o.Value}
	case *Integer:
		return &Float{Value: self.Value * float64(o.Value)}
	}
	return nil
}
func (self *Float) Div(other BaseValueInterface) BaseValueInterface {
	switch o := other.(type) {
	case *Float:
		return &Float{Value: self.Value / o.Value}
	case *Integer:
		return &Float{Value: self.Value / float64(o.Value)}
	}
	return nil
}

func (self *Float) String() string {
	return strconv.FormatFloat(self.Value, 'g', -1, 64)
}
