package errors

import "fmt"

type Error struct {
	Name    string
	Details string
}

func (e Error) String() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Details)
}

func NewIllegalCharError(details string) *Error {
	return &Error{Name: "Illegal Character", Details: details}
}

func NewInvalidNumberError(details string) *Error {
	return &Error{Name: "Invalid Number", Details: details}
}
