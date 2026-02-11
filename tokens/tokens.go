package tokens

import "fmt"

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
	Type  TokenType
	Value any
}

func (t Token) String() string {
	if t.Value == nil {
		return fmt.Sprintf("Token{Type=%s}", TokenTypeName[t.Type])
	}
	return fmt.Sprintf("Token{Type=%s, Value=%v}", TokenTypeName[t.Type], t.Value)
}
