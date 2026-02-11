package lexer

import (
	"strconv"
	"unicode"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/position"
	"github.com/DGTV11/weh-script/tokens"
)

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

func (l *Lexer) MakeNumberToken() (*tokens.Token, error) {
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
		_type tokens.TokenType
	)

	if dotCount == 0 {
		value, err = strconv.ParseInt(numStr, 10, 64)
		_type = tokens.TokenTypeInt
	} else {
		value, err = strconv.ParseFloat(numStr, 64)
		_type = tokens.TokenTypeFloat
	}

	if err != nil {
		return nil, err
	}

	return &tokens.Token{Type: _type, Value: value}, nil
}

func (l *Lexer) Tokenise() ([]tokens.Token, *errors.Error) {
	var tokenList []tokens.Token

	for l.CurrentChar != nil {
		switch char := *l.CurrentChar; char {
		case ' ':
		case '\t':
		case '+':
			tokenList = append(tokenList, tokens.Token{Type: tokens.TokenTypePlus, Value: nil})
		case '-':
			tokenList = append(tokenList, tokens.Token{Type: tokens.TokenTypeMinus, Value: nil})
		case '*':
			tokenList = append(tokenList, tokens.Token{Type: tokens.TokenTypeMul, Value: nil})
		case '/':
			tokenList = append(tokenList, tokens.Token{Type: tokens.TokenTypeDiv, Value: nil})
		case '(':
			tokenList = append(tokenList, tokens.Token{Type: tokens.TokenTypeLparen, Value: nil})
		case ')':
			tokenList = append(tokenList, tokens.Token{Type: tokens.TokenTypeRparen, Value: nil})
		default:
			if unicode.IsDigit(char) {
				tokp, err := l.MakeNumberToken()
				if err != nil {
					positionStart := l.Position.Copy()
					l.Advance()
					return []tokens.Token{}, errors.NewInvalidNumberError(positionStart, l.Position, err.Error())
				}
				tokenList = append(tokenList, *tokp)

				continue
			}

			positionStart := l.Position.Copy()
			l.Advance()
			return []tokens.Token{}, errors.NewIllegalCharError(positionStart, l.Position, "'"+string(char)+"'")
		}

		l.Advance()
	}

	return tokenList, nil
}
