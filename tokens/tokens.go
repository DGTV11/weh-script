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
	TokenTypeEOF
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
	TokenTypeEOF:    "TokenTypeEOF",
}

type Token struct {
	Type     TokenType
	Value    any
	PosStart *position.Position
	PosEnd   *position.Position
}

func NewToken(_type TokenType, value any, posStartIn *position.Position, posEndIn *position.Position) Token {
	var posStart *position.Position
	var posEnd *position.Position

	if posStartIn != nil {
		posStart = posStartIn.Copy()
		posEnd = posStartIn.Copy()
		posEnd.Advance(nil)
	}

	if posEndIn != nil {
		posEnd = posEndIn.Copy()
	}

	return Token{Type: _type, Value: value, PosStart: posStart, PosEnd: posEnd}
}

func (t Token) String() string {
	if t.Value == nil {
		// return fmt.Sprintf("Token{Type=%s}", TokenTypeName[t.Type])
		return fmt.Sprintf("%s", TokenTypeName[t.Type])
	}
	// return fmt.Sprintf("Token{Type=%s, Value=%v}", TokenTypeName[t.Type], t.Value)
	return fmt.Sprintf("%s:%v", TokenTypeName[t.Type], t.Value)
}
