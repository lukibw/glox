package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lukibw/glox/interpreter"
	"github.com/lukibw/glox/parser"
	"github.com/lukibw/glox/scanner"
)

func main() {
	source, err := os.ReadFile("main.lox")
	if err != nil {
		log.Fatalln(err)
	}
	scanner := scanner.New(string(source))
	tokens, err := scanner.Run()
	if err != nil {
		log.Fatalln(err)
	}
	parser := parser.New(tokens)
	statements, errs := parser.Run()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	interpreter := interpreter.New()
	if err = interpreter.Run(statements); err != nil {
		log.Fatalln(err)
	}
}
