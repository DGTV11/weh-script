package values

import (
	"github.com/DGTV11/weh-script/position"
)

type Value interface {
	SetPos(positionStart position.Position, positionEnd position.Position) //TODO: need or not?
	Add(first Value, second Value) Value
	Sub(first Value, second Value) Value
	Mul(first Value, second Value) Value
	Div(first Value, second Value) Value
	Repr(first Value, second Value) Value //or String() ?
}

//TODO: https://youtu.be/YYvBy0vqcSw?si=tZCzn9lsaC2v8dPf&t=363
//*need to add Number value member
