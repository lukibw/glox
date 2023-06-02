package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lukibw/glox/expr"
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
	parser := parser.New(tokens)
	exp, err := parser.Parse()
	if err != nil {
		log.Fatalln(err)
	}
	printer := expr.NewPrinter()
	fmt.Println(printer.Print(exp))
}
