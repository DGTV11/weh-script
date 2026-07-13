package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
	// "reflect"

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
	if viewExecutionTimes == true {
		tmpStartTime = time.Now()
		tmpTime = time.Now()
	}
	_lexer := lexer.NewLexer(fileName, text)
	tokens, err := _lexer.Tokenise()
	if viewExecutionTimes == true {
		lexerElapsed = time.Now().Sub(tmpTime)
	}
	if err != nil {
		return nil, err
	}

	if viewExecutionTimes == true {
		tmpTime = time.Now()
	}
	_parser := parser.NewParser(tokens)
	ast := _parser.Parse()
	if viewExecutionTimes == true {
		parserElapsed = time.Now().Sub(tmpTime)
	}
	if ast.Err != nil {
		return nil, ast.Err
	}
	// fmt.Println(ast.Node)

	if viewExecutionTimes == true {
		tmpTime = time.Now()
	}
	context := environment.Context{DisplayName: "<program>", SymTable: globalSymbolTable}
	result := interpreter.Visit(ast.Node, context)
	if viewExecutionTimes == true {
		interpreterElapsed = time.Now().Sub(tmpTime)
		elapsed = time.Now().Sub(tmpStartTime)
	}

	return result.Value, result.Err
}

var text string
var bytecodeMode bool
var viewExecutionTimes bool
var tmpTime time.Time
var tmpStartTime time.Time
var lexerElapsed time.Duration
var parserElapsed time.Duration
var interpreterElapsed time.Duration
var elapsed time.Duration

func shell() {
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
			}
			// if reflect.TypeOf(res).String() != "*values.Null" {
			if _, ok := res.(*values.Null); ok == false {
				fmt.Println(res)

				globalSymbolTable.SetSymbol("_", res)
			}
			if viewExecutionTimes == true {
				fmt.Printf("\n===========================\nlexer %v\nparser %v\ninterpreter %v\n---------------------------\ntotal %v\n===========================\n", lexerElapsed, parserElapsed, interpreterElapsed, elapsed)
			}
		}
	}
}

func runFile(fp string) {
	if fileExt := filepath.Ext(fp); fileExt != ".weh" {
		log.Fatal("Invalid file extension") //TODO: .wvm bytecode files
	}

	programBytestr, Rerr := os.ReadFile(fp)
	if Rerr != nil {
		log.Fatal(Rerr)
	}

	program := string(programBytestr)
	globalSymbolTable := SetupGlobalSymbolTable()
	_, err := Run(fp, program, globalSymbolTable)

	if err != nil {
		fmt.Println(err)
	}
	if viewExecutionTimes == true {
		fmt.Printf("\n===========================\nlexer %v\nparser %v\ninterpreter %v\n---------------------------\ntotal %v\n===========================\n", lexerElapsed, parserElapsed, interpreterElapsed, elapsed)
	}

}

func main() {
	pflag.BoolVar(&bytecodeMode, "bytecode-mode", false, "Enable bytecode mode") //after implementing bytecode VM: default this to true and make non-bytecode mode legacy
	pflag.BoolVar(&viewExecutionTimes, "time", false, "View execution times")
	pflag.Parse()

	if fp := pflag.Arg(0); fp == "" {
		shell()
	} else {
		runFile(fp)
	}
}
