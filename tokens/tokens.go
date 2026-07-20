package tokens

import (
	"fmt"

	"github.com/DGTV11/weh-script/position"
)

type TokenType int

const (
	TokenTypeInt TokenType = iota
	TokenTypeFloat
	TokenTypeString
	TokenTypeChar
	TokenTypeIdentifier
	TokenTypeKeyword
	TokenTypePlus
	TokenTypeMinus
	TokenTypeMul
	TokenTypeDiv
	TokenTypePow
	TokenTypeEquals
	TokenTypeLparen
	TokenTypeRparen
	TokenTypeLsquare
	TokenTypeRsquare
	TokenTypeEE
	TokenTypeNE
	TokenTypeLT
	TokenTypeGT
	TokenTypeLTE
	TokenTypeGTE
	TokenTypeLAnd
	TokenTypeLOr
	TokenTypeLNot
	TokenTypeBAnd
	TokenTypeBOr
	TokenTypeBNot
	TokenTypeComma
	TokenTypeArrow
	TokenTypeNewline
	TokenTypeEOF
)

var TokenTypeNameMap = [...]string{
	TokenTypeInt:        "TokenTypeInt",
	TokenTypeFloat:      "TokenTypeFloat",
	TokenTypeString:     "TokenTypeString",
	TokenTypeChar:       "TokenTypeChar",
	TokenTypeIdentifier: "TokenTypeIdentifier",
	TokenTypeKeyword:    "TokenTypeKeyword",
	TokenTypePlus:       "TokenTypePlus",
	TokenTypeMinus:      "TokenTypeMinus",
	TokenTypeMul:        "TokenTypeMul",
	TokenTypeDiv:        "TokenTypeDiv",
	TokenTypePow:        "TokenTypePow",
	TokenTypeEquals:     "TokenTypeEquals",
	TokenTypeLparen:     "TokenTypeLparen",
	TokenTypeRparen:     "TokenTypeRparen",
	TokenTypeLsquare:    "TokenTypeLsquare",
	TokenTypeRsquare:    "TokenTypeRsquare",
	TokenTypeEE:         "TokenTypeEE",
	TokenTypeNE:         "TokenTypeNE",
	TokenTypeLT:         "TokenTypeLT",
	TokenTypeGT:         "TokenTypeGT",
	TokenTypeLTE:        "TokenTypeLTE",
	TokenTypeGTE:        "TokenTypeGTE",
	TokenTypeLAnd:       "TokenTypeLAnd",
	TokenTypeLOr:        "TokenTypeLOr",
	TokenTypeLNot:       "TokenTypeLNot",
	TokenTypeBAnd:       "TokenTypeBAnd",
	TokenTypeBOr:        "TokenTypeBOr",
	TokenTypeBNot:       "TokenTypeBNot",
	TokenTypeComma:      "TokenTypeComma",
	TokenTypeArrow:      "TokenTypeArrow",
	TokenTypeNewline:    "TokenTypeNewline",
	TokenTypeEOF:        "TokenTypeEOF",
}

var Keywords = []string{
	"var",
	"nonlocal",
	"del",
	"if",
	"elif",
	"else",
	"then",
	"for",
	"to",
	"step",
	"while",
	"func",
	"struct",
	"end",
	"return",
	"continue",
	"break",
	"import",
}

type Token struct {
	Type     TokenType
	Value    any
	PosRange position.PositionRange
}

type TokenTV struct {
	Type  TokenType
	Value any
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

	return Token{Type: _type, Value: value, PosRange: position.PositionRange{Start: posStart, End: posEnd}}
}

func (t Token) Matches(type_ TokenType, value any) bool {
	return t.Type == type_ && t.Value == value
}

func (t Token) String() string {
	if t.Value == nil {
		return fmt.Sprintf("%s", TokenTypeNameMap[t.Type])
	}
	return fmt.Sprintf("%s:%v", TokenTypeNameMap[t.Type], t.Value)
}
