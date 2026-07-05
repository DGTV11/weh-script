package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/interpreter"
	"github.com/DGTV11/weh-script/lexer"
	"github.com/DGTV11/weh-script/parser"
	"github.com/DGTV11/weh-script/runtime"
	"github.com/DGTV11/weh-script/values"
)

func SetupGlobalymbolTable() *runtime.SymbolTable {
	GlobalSymbolTable := runtime.SymbolTable{Symbols: map[string]any{}}

	GlobalSymbolTable.SetSymbol("null", &values.Integer{Value: 0})
	GlobalSymbolTable.SetSymbol("true", &values.Integer{Value: 1})
	GlobalSymbolTable.SetSymbol("false", &values.Integer{Value: 0})

	return &GlobalSymbolTable
}

func Run(fileName string, text string, globalSymbolTable *runtime.SymbolTable) (any, *errors.Error) {
	_lexer := lexer.NewLexer(fileName, text)
	tokens, err := _lexer.Tokenise()
	if err != nil {
		return nil, err
	}

	_parser := parser.NewParser(tokens)
	ast := _parser.Parse()
	if ast.Err != nil {
		return nil, ast.Err
	}

	context := runtime.Context{DisplayName: "<program>", SymTable: globalSymbolTable}
	result := interpreter.Visit(ast.Node, context)

	return result.Value, result.Err
}

var text string

func main() {
	fmt.Println("WehScript Programming Language")

	globalSymbolTable := SetupGlobalymbolTable()

	for true {
		fmt.Print("wehscript > ")
		// fmt.Scanf("%[^\n]", &text)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		serr := scanner.Err()
		if serr != nil {
			fmt.Println(serr)
		}

		res, err := Run("<stdin>", scanner.Text(), globalSymbolTable)

		if err != nil {
			fmt.Println(err)
		} else if res != nil {
			fmt.Println(res)
		}
	}
}
