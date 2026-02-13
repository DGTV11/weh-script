package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/nodes"

	"github.com/DGTV11/weh-script/lexer"
	"github.com/DGTV11/weh-script/parser"
)

func run(fileName string, text string) (nodes.Node, *errors.Error) {
	_lexer := lexer.NewLexer(fileName, text)
	tokens, err := _lexer.Tokenise()
	if err != nil {
		return nil, err
	}

	_parser := parser.NewParser(tokens)
	ast := _parser.Parse()

	return ast.Node, ast.Err
}

var text string

func main() {
	fmt.Println("WehScript Programming Language")

	for true {
		fmt.Print("wehscript > ")
		// fmt.Scanf("%[^\n]", &text)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		serr := scanner.Err()
		if serr != nil {
			fmt.Println(serr)
		}

		res, err := run("<stdin>", scanner.Text())

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	}
}
