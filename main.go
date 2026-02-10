package main

import (
	"fmt"
	"github.com/DGTV11/weh-script/lexer"
)

var running bool = true
var text string

func main() {
	fmt.Println("WehScript Programming Language")
	for running {
		fmt.Print("wehscript > ")
		fmt.Scanln(&text)
	}
}
