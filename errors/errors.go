package errors

import (
	"fmt"

	"github.com/DGTV11/weh-script/position"
)

type Error struct {
	PositionStart position.Position
	PositionEnd   position.Position
	Name          string
	Details       string
}

func (e Error) String() string {
	return fmt.Sprintf("%s: %s\nFile %s, line %d", e.Name, e.Details, e.PositionStart.FileName, e.PositionStart.Line+1)
}

func NewIllegalCharError(positionStart position.Position, positionEnd position.Position, details string) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Illegal Character", Details: details}
}

func NewInvalidNumberError(positionStart position.Position, positionEnd position.Position, details string) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Invalid Number", Details: details}
}
