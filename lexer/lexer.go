package lexer

import (
	"fmt"
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

type Token struct {
	_type TokenType
	value any //TODO
}

type Lexer struct {
	Text        string
	Pos         int
	CurrentChar *rune
}

func (l *Lexer) Advance() {
	l.Pos += 1
	if l.Pos >= len(l.Text) {
		l.CurrentChar = nil
		return
	}
	l.CurrentChar = &[]rune(l.Text)[l.Pos]
}

func (l *Lexer) Tokenise() *Token {
	var tokens []Token

	for l.CurrentChar != nil {
		if *l.CurrentChar == ' ' || *l.CurrentChar == '\t' {
			l.Advance()
		}
	}
	tokens = append(tokens, nil) //TODO
}

func NewLexer(input string) *Lexer {
	newLexer := Lexer{Text: input, Pos: -1, CurrentChar: nil}
	newLexer.Advance()
	return &newLexer
}
