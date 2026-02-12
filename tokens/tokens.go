package tokens

import (
	"fmt"
	"github.com/DGTV11/weh-script/position"
)

type TokenType int

const (
	TokenTypeInt TokenType = iota
	TokenTypeFloat
	TokenTypePlus
	TokenTypeMinus
	TokenTypeMul
	TokenTypeDiv
	TokenTypeLparen
	TokenTypeRparen
)

var TokenTypeName = map[TokenType]string{
	TokenTypeInt:    "TokenTypeInt",
	TokenTypeFloat:  "TokenTypeFloat",
	TokenTypePlus:   "TokenTypePlus",
	TokenTypeMinus:  "TokenTypeMinus",
	TokenTypeMul:    "TokenTypeMul",
	TokenTypeDiv:    "TokenTypeDiv",
	TokenTypeLparen: "TokenTypeLparen",
	TokenTypeRparen: "TokenTypeRparen",
}

type Token struct {
	Type     TokenType
	Value    any
	PosStart position.Position
	PosEnd   position.Position
}

func NewToken(_type TokenType, value any, posStart *position.Position, posEnd *position.Position) Token {
	//TODO: implement NewToken which determines wtf to put in Position (also update lexer accordingly)
	return Token{Type: _type, Value: value, PosStart: todo, PosEnd: todo} //TODO
}

func (t Token) String() string {
	if t.Value == nil {
		return fmt.Sprintf("Token{Type=%s}", TokenTypeName[t.Type])
	}
	return fmt.Sprintf("Token{Type=%s, Value=%v}", TokenTypeName[t.Type], t.Value) //TODO: add pos
}
