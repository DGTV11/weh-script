package main

import (
	"fmt"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/lexer"
)

func run(text string) ([]lexer.Token, *errors.Error) {
	lexer := lexer.NewLexer(text)
	tokens, err := lexer.Tokenise()
	return tokens, err
}

var running bool = true
var text string

func main() {
	fmt.Println("WehScript Programming Language")
	for running {
		fmt.Print("wehscript > ")
		fmt.Scanln(&text)

		res, err := run(text)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	}
}
