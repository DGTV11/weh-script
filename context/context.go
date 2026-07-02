package context

import (
	"fmt"

	"github.com/DGTV11/weh-script/position"
)

type Context struct {
	DisplayName    string
	Parent         *Context
	ParentEntryPos *position.Position
}

func (ctx *Context) GenerateTraceback(positionStart *position.Position) string {
	result := ""
	pos := positionStart
	currentCtx := ctx

	for currentCtx != nil {
		result += fmt.Sprintf("\tFile %s, line %d, in %s\n", pos.FileName, pos.Line+1, currentCtx.DisplayName)
		pos = currentCtx.ParentEntryPos
		currentCtx = currentCtx.Parent
	}
	return "Traceback (most recent call last):\n" + result
}
