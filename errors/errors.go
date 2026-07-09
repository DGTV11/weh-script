package errors

import (
	"fmt"
	"strings"

	"github.com/DGTV11/weh-script/environment"
	"github.com/DGTV11/weh-script/position"
)

// *Regular errors
func StringWithArrows(text string, positionStart *position.Position, positionEnd *position.Position) string {
	result := ""

	// calculate indices
	indexStart := max(strings.LastIndexByte(text[:positionStart.Index], '\n'), 0)
	indexEnd := strings.IndexByte(text[:indexStart], '\n')
	if indexEnd < 0 {
		indexEnd = len(text)
	}

	// generate each line
	lineCount := positionEnd.Line - positionStart.Line + 1
	for i := 0; i < lineCount; i++ {
		// calculate line cols
		line := text[indexStart:indexEnd]
		columnStart := 0
		columnEnd := 0
		if i == 0 {
			columnStart = positionStart.Column
		} else {
			columnStart = 0
		}
		if i == lineCount-1 {
			columnEnd = positionEnd.Column
		} else {
			columnEnd = len(line) - 1
		}

		// append to result
		result += line + "\n"
		result += strings.Repeat(" ", columnStart) + strings.Repeat("^", (columnEnd-columnStart))

		// recalculate indices
		indexStart = indexEnd
		indexEnd = strings.IndexByte(text[:indexStart], '\n')
		if indexEnd < 0 {
			indexEnd = len(text)
		}
	}

	return strings.ReplaceAll(result, "\t", "")
}

type Error struct {
	PositionStart *position.Position
	PositionEnd   *position.Position
	Name          string
	Details       string
	Ctx           *environment.Context
}

func (e Error) String() string {
	errString := fmt.Sprintf("%s: %s\nFile %s, line %d\n\n%s", e.Name, e.Details, e.PositionStart.FileName, e.PositionStart.Line+1, StringWithArrows(e.PositionStart.FileText, e.PositionStart, e.PositionEnd))

	if e.Ctx == nil {
		return errString
	}
	return e.Ctx.GenerateTraceback(e.PositionStart) + errString
}

func NewIllegalCharError(positionStart *position.Position, positionEnd *position.Position, details string) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Illegal Character", Details: details}
}

func NewInvalidNumberError(positionStart *position.Position, positionEnd *position.Position, details string) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Invalid Number", Details: details}
}

func NewExpectedCharError(positionStart *position.Position, positionEnd *position.Position, details string) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Expected Character", Details: details}
}

func NewInvalidSyntaxError(positionStart *position.Position, positionEnd *position.Position, details string) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Invalid Syntax", Details: details}
}

func NewSyntaxNotImplementedError(positionStart *position.Position, positionEnd *position.Position, details string) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Syntax Not Implemented", Details: details}
}

//*Runtime errors

func NewRuntimeError(positionStart *position.Position, positionEnd *position.Position, details string, ctx environment.Context) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Runtime Error", Details: details, Ctx: &ctx}
}

func NewNotImplementedError(positionStart *position.Position, positionEnd *position.Position, details string, ctx environment.Context) *Error {
	return &Error{PositionStart: positionStart, PositionEnd: positionEnd, Name: "Not Implemented", Details: details, Ctx: &ctx}
}
