package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	// "reflect"
	"strings"
	"unsafe"

	"github.com/spf13/pflag"

	"github.com/DGTV11/weh-script/environment"
	"github.com/DGTV11/weh-script/errors"
	"github.com/DGTV11/weh-script/interpreter"
	"github.com/DGTV11/weh-script/lexer"
	"github.com/DGTV11/weh-script/parser"
	"github.com/DGTV11/weh-script/values"
)

const _ = uint(1) / (uint(unsafe.Sizeof(int(0))) - 7) //ensures that size of int == size of int64

func SetupGlobalSymbolTable() *environment.SymbolTable {
	GlobalSymbolTable := environment.SymbolTable{Symbols: map[string]any{}}

	//*Load constants
	GlobalSymbolTable.SetSymbol("null", &values.Null{})
	GlobalSymbolTable.SetSymbol("true", &values.Integer{Value: 1})
	GlobalSymbolTable.SetSymbol("false", &values.Integer{Value: 0})

	//*Load functions
	for funcName := range maps.Keys(interpreter.BuiltInFunctionTable) {
		GlobalSymbolTable.SetSymbol(funcName, &values.BuiltInFunction{BaseFunction: values.BaseFunction{Name: &funcName}})
	}

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
var bytecodeMode bool

func main() {

	pflag.BoolVar(&bytecodeMode, "bytecode-mode", false, "Enable bytecode mode") //after implementing bytecode VM: default this to true and make non-bytecode mode legacy
	pflag.Parse()

	fmt.Println("WehScript Programming Language")
	if bytecodeMode == true {
		log.Fatal("Bytecode VM not implemented")
	} else {
		globalSymbolTable := SetupGlobalSymbolTable()

		for {
			fmt.Print("wehscript > ")
			// fmt.Scanf("%[^\n]", &text)

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			serr := scanner.Err()
			if serr != nil {
				fmt.Println(serr)
			}

			text := scanner.Text()
			if strings.TrimSpace(text) == "" {
				continue
			}

			res, err := Run("<stdin>", text, globalSymbolTable)

			if err != nil {
				fmt.Println(err)
				// } else if reflect.TypeOf(res).String() != "*values.Null" {
			} else if res != nil {
				resV := res.(*values.List)
				if len(resV.Elements) == 1 {
					fmt.Println(resV.Elements[0])
				} else {
					fmt.Println(res)
				}

				// globalSymbolTable.SetSymbol("_", res) //TODO: update '_' variable after every expression (separate statementsnode?)
			}
		}
	}
}
