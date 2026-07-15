package environment

import (
	"fmt"

	"github.com/DGTV11/weh-script/position"
)

// *Context
type Context struct {
	DisplayName    string
	Parent         *Context
	ParentEntryPos *position.Position
	SymTable       *SymbolTable
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

// *SymbolTable
type SymbolTable struct {
	Symbols map[string]any
	Parent  *SymbolTable
}

func (s SymbolTable) GetSymbol(name string) any {
	val, ok := s.Symbols[name]
	if ok == true {
		return val
	}
	if s.Parent != nil {
		return s.Parent.GetSymbol(name)
	}
	return nil
}

func (s SymbolTable) ForceSetSymbol(name string, val any) {
	s.Symbols[name] = val
}

func (s SymbolTable) SetSymbol(name string, val any) bool {
	_, ok := s.Symbols[name]
	if ok == false {
		s.Symbols[name] = val
	}
	return !ok
}

func (s SymbolTable) UpdateSymbol(name string, val any) bool {
	_, ok := s.Symbols[name]
	if ok == true {
		s.Symbols[name] = val
	}
	return ok
}

func (s SymbolTable) RemoveSymbol(name string) {
	delete(s.Symbols, name)
}
