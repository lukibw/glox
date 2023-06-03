package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lukibw/glox/ast"
	"github.com/lukibw/glox/parser"
	"github.com/lukibw/glox/scanner"
)

func main() {
	source, err := os.ReadFile("main.lox")
	if err != nil {
		log.Fatalln(err)
	}
	scanner := scanner.New(string(source))
	tokens, err := scanner.ScanTokens()
	if err != nil {
		log.Fatalln(err)
	}
	parser := parser.New[any](tokens)
	expr, err := parser.Parse()
	if err != nil {
		log.Fatalln(err)
	}
	interpreter := ast.NewInterpreter()
	result, err := interpreter.Interpret(expr)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}
