package lexer

import (
	"slices"
	"strconv"
	"unicode"

	"github.com/stanNthe5/stringbuf"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/position"
	"github.com/DGTV11/weh-script/tokens"
)

type Lexer struct {
	FileName    string
	Text        []rune
	Position    position.Position
	CurrentChar *rune
}

func NewLexer(fileName string, text string) *Lexer {
	newLexer := Lexer{
		FileName:    fileName,
		Text:        []rune(text),
		Position:    *position.NewPosition(-1, 0, -1, fileName, text),
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
	l.CurrentChar = &l.Text[l.Position.Index]
}

func (l *Lexer) Tokenise() ([]tokens.Token, *errors.Error) {
	var tokenList []tokens.Token

	for l.CurrentChar != nil {
		switch char := *l.CurrentChar; char {
		case ' ':
			l.Advance()
		case '\t':
			l.Advance()
		case '#':
			l.Advance()
			for l.CurrentChar != nil && *l.CurrentChar != '\n' && *l.CurrentChar != ';' {
				// for l.CurrentChar != nil && *l.CurrentChar != '\n' {
				l.Advance()
			}
			if l.CurrentChar != nil {
				l.Advance()
			}
		case '\n':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeNewline, nil, &l.Position, nil))
			l.Advance()
		case ';':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeNewline, nil, &l.Position, nil))
			l.Advance()
		case '"':
			posStart := l.Position.Copy()
			tokp := l.MakeString()
			if tokp == nil {
				return []tokens.Token{}, errors.NewInvalidSyntaxError(posStart, &l.Position, "Unterminated string")
			}
			tokenList = append(tokenList, *tokp)
		case '+':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypePlus, nil, &l.Position, nil))
			l.Advance()
		case '-':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeMinus, nil, &l.Position, nil))
			l.Advance()
		case '*':
			posStart := l.Position.Copy()
			l.Advance()
			if l.CurrentChar != nil && *l.CurrentChar == '*' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypePow, nil, posStart, &l.Position))
			} else {
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeMul, nil, posStart, nil))
			}
		case '/':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeDiv, nil, &l.Position, nil))
			l.Advance()
		case '!':
			posStart := l.Position.Copy()
			l.Advance()
			if l.CurrentChar != nil && *l.CurrentChar == '=' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeNE, nil, posStart, &l.Position))
			} else {
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeLNot, nil, posStart, nil))
			}
		case '=':
			posStart := l.Position.Copy()
			l.Advance()
			if l.CurrentChar != nil && *l.CurrentChar == '=' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeEE, nil, posStart, &l.Position))
			} else if l.CurrentChar != nil && *l.CurrentChar == '>' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeArrow, nil, posStart, &l.Position))
			} else {
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeEquals, nil, posStart, nil))
			}
		case '<':
			posStart := l.Position.Copy()
			l.Advance()
			if l.CurrentChar != nil && *l.CurrentChar == '=' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeLTE, nil, posStart, &l.Position))
			} else {
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeLT, nil, posStart, nil))
			}
		case '>':
			posStart := l.Position.Copy()
			l.Advance()
			if l.CurrentChar != nil && *l.CurrentChar == '=' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeGTE, nil, posStart, &l.Position))
			} else {
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeGT, nil, posStart, nil))
			}
		case '&':
			posStart := l.Position.Copy()
			l.Advance()
			if l.CurrentChar != nil && *l.CurrentChar == '&' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeLAnd, nil, posStart, &l.Position))
			} else {
				// tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeBAnd, nil, posStart, nil))
				return []tokens.Token{}, errors.NewSyntaxNotImplementedError(posStart, &l.Position, "Bitwise 'and' not implemented")
			}
		case '|':
			posStart := l.Position.Copy()
			l.Advance()
			if l.CurrentChar != nil && *l.CurrentChar == '|' {
				l.Advance()
				tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeLOr, nil, posStart, &l.Position))
			} else {
				// tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeBOr, nil, posStart, nil))
				return []tokens.Token{}, errors.NewSyntaxNotImplementedError(posStart, &l.Position, "Bitwise 'or' not implemented")
			}
		case '~':
			// tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeBNot, nil, &l.Position, nil))
			// l.Advance()
			posStart := l.Position.Copy()
			l.Advance()
			return []tokens.Token{}, errors.NewSyntaxNotImplementedError(posStart, &l.Position, "Bitwise 'not' not implemented")
		case '(':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeLparen, nil, &l.Position, nil))
			l.Advance()
		case ')':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeRparen, nil, &l.Position, nil))
			l.Advance()
		case '[':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeLsquare, nil, &l.Position, nil))
			l.Advance()
		case ']':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeRsquare, nil, &l.Position, nil))
			l.Advance()
		case ',':
			tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeComma, nil, &l.Position, nil))
			l.Advance()
		default:
			if unicode.IsDigit(char) {
				tokp, err := l.MakeNumberToken()
				if err != nil {
					positionStart := l.Position.Copy()
					l.Advance()
					return []tokens.Token{}, errors.NewInvalidNumberError(positionStart, &l.Position, err.Error())
				}
				tokenList = append(tokenList, *tokp)

				continue
			} else if unicode.IsLetter(char) || char == '_' {
				tokp := l.MakeIdentifierOrKeywordToken()
				tokenList = append(tokenList, *tokp)

				continue
			}

			positionStart := l.Position.Copy()
			l.Advance()
			return []tokens.Token{}, errors.NewIllegalCharError(positionStart, &l.Position, "'"+string(char)+"'")
		}

	}

	tokenList = append(tokenList, tokens.NewToken(tokens.TokenTypeEOF, nil, &l.Position, nil))

	return tokenList, nil
}

func (l *Lexer) MakeNumberToken() (*tokens.Token, error) {
	numStr := ""
	dotCount := 0
	posStart := l.Position.Copy()

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
			dotCount++
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

	// return &tokens.Token{Type: _type, Value: value}, nil
	newTok := tokens.NewToken(_type, value, posStart, &l.Position)
	return &newTok, nil
}

func (l *Lexer) MakeIdentifierOrKeywordToken() *tokens.Token {
	idStr := ""
	posStart := l.Position.Copy()

	for l.CurrentChar != nil {
		char := *l.CurrentChar

		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' {
			idStr += string(char)

			l.Advance()

			continue
		}

		break
	}

	var _type tokens.TokenType

	if slices.Contains(tokens.Keywords, idStr) {
		_type = tokens.TokenTypeKeyword
	} else {
		_type = tokens.TokenTypeIdentifier
	}

	newTok := tokens.NewToken(_type, idStr, posStart, &l.Position)
	return &newTok
}

var EscapeChars = map[rune]rune{
	'n': '\n',
	't': '\t',
	'r': '\r',
}

func (l *Lexer) MakeString() *tokens.Token {
	sb := stringbuf.New("")
	posStart := l.Position.Copy()
	atEscapeChar := false
	l.Advance()

	for l.CurrentChar != nil && (*l.CurrentChar != '"' || atEscapeChar == true) {
		if atEscapeChar {
			escapedRune, ok := EscapeChars[*l.CurrentChar]
			if ok == false {
				escapedRune = *l.CurrentChar
			}
			sb.AppendRune(escapedRune)
			atEscapeChar = false
		} else {
			if *l.CurrentChar == '\\' {
				atEscapeChar = true
			} else {
				sb.AppendRune(*l.CurrentChar)
				atEscapeChar = false
			}
		}

		l.Advance()
	}

	if l.CurrentChar == nil || *l.CurrentChar != '"' {
		return nil
	}

	l.Advance()
	newTok := tokens.NewToken(tokens.TokenTypeString, sb.String(), posStart, &l.Position)
	return &newTok
}
