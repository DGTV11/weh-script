package lexer

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/DGTV11/weh-script/errors"
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
	Type  TokenType
	Value any
}

func (t Token) String() string {
	if t.Value == nil {
		return fmt.Sprintf("Token{Type=%s}", TokenTypeName[t.Type])
	}
	return fmt.Sprintf("Token{Type=%s, Value=%v}", TokenTypeName[t.Type], t.Value)
}

type Lexer struct {
	FileName    string
	Text        string
	Position    position.Position
	CurrentChar *rune
}

func NewLexer(fileName string, text string) *Lexer {
	newLexer := Lexer{
		FileName:    fileName,
		Text:        text,
		Position:    position.NewPosition(-1, 0, -1, fileName, text),
		CurrentChar: nil,
	}
	newLexer.Advance()
	return &newLexer
}

func (l *Lexer) Advance() {
	l.Position.Advance(l.CurrentChar)
	if l.Position.Index >= len(l.Text) {
		l.CurrentChar = nil
		return
	}
	l.CurrentChar = &[]rune(l.Text)[l.Position.Index]
}

func (l *Lexer) MakeNumberToken() (*Token, error) {
	numStr := ""
	dotCount := 0

	for l.CurrentChar != nil {
		char := *l.CurrentChar

		if unicode.IsDigit(char) {
			numStr += string(char)

			l.Advance()

			continue
		}

		if char == '.' {
			if dotCount == 1 {
				break
			}
			dotCount += 1
			numStr += "."

			l.Advance()

			continue
		}

		break
	}

	var (
		value any
		err   error
		_type TokenType
	)

	if dotCount == 0 {
		value, err = strconv.ParseInt(numStr, 10, 64)
		_type = TokenTypeInt
	} else {
		value, err = strconv.ParseFloat(numStr, 64)
		_type = TokenTypeFloat
	}

	if err != nil {
		return nil, err
	}

	return &Token{Type: _type, Value: value}, nil
}

func (l *Lexer) Tokenise() ([]Token, *errors.Error) {
	var tokens []Token

	for l.CurrentChar != nil {
		switch char := *l.CurrentChar; char {
		case ' ':
		case '\t':
		case '+':
			tokens = append(tokens, Token{Type: TokenTypePlus, Value: nil})
		case '-':
			tokens = append(tokens, Token{Type: TokenTypeMinus, Value: nil})
		case '*':
			tokens = append(tokens, Token{Type: TokenTypeMul, Value: nil})
		case '/':
			tokens = append(tokens, Token{Type: TokenTypeDiv, Value: nil})
		case '(':
			tokens = append(tokens, Token{Type: TokenTypeLparen, Value: nil})
		case ')':
			tokens = append(tokens, Token{Type: TokenTypeRparen, Value: nil})
		default:
			if unicode.IsDigit(char) {
				tokp, err := l.MakeNumberToken()
				if err != nil {
					positionStart := l.Position.Copy()
					l.Advance()
					return []Token{}, errors.NewInvalidNumberError(positionStart, l.Position, err.Error())
				}
				tokens = append(tokens, *tokp)

				continue
			}

			positionStart := l.Position.Copy()
			l.Advance()
			return []Token{}, errors.NewIllegalCharError(positionStart, l.Position, "'"+string(char)+"'")
		}

		l.Advance()
	}

	return tokens, nil
}
