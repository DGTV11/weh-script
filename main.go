package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/DGTV11/weh-script/compiler/environment"
	"github.com/DGTV11/weh-script/compiler/errors"
	"github.com/DGTV11/weh-script/compiler/interpreter"
	"github.com/DGTV11/weh-script/compiler/lexer"
	"github.com/DGTV11/weh-script/compiler/parser"
	"github.com/DGTV11/weh-script/compiler/values"
)

const _ = uint(1) / (uint(unsafe.Sizeof(int(0))) - 7) //ensures that size of int == size of int64

func SetupGlobalymbolTable() *environment.SymbolTable {
	GlobalSymbolTable := environment.SymbolTable{Symbols: map[string]any{}}

	GlobalSymbolTable.SetSymbol("null", &values.Null{})
	GlobalSymbolTable.SetSymbol("true", &values.Integer{Value: 1})
	GlobalSymbolTable.SetSymbol("false", &values.Integer{Value: 0})

	return &GlobalSymbolTable
}

func Run(fileName string, text string, globalSymbolTable *environment.SymbolTable) (any, *errors.Error) {
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
	// fmt.Println(ast.Node)

	context := environment.Context{DisplayName: "<program>", SymTable: globalSymbolTable}
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
		} else if reflect.TypeOf(res).String() != "*values.Null" {
			fmt.Println(res)
		}
	}
}
