package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lukibw/glox/scan"
)

func main() {
	source, err := os.ReadFile("main.lox")
	if err != nil {
		log.Fatalln(err)
	}
	lexer := scan.NewLexer(string(source))
	tokens, err := lexer.ScanTokens()
	if err != nil {
		log.Fatalln(err)
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
}
