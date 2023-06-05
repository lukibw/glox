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
	statements, errs := parser.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	interpreter := ast.NewInterpreter()
	err = interpreter.Interpret(statements)
	if err != nil {
		log.Fatalln(err)
	}
}
