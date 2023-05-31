package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lukibw/glox/ast"
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
	binaryExpr := &ast.BinaryExpr{
		Left: &ast.UnaryExpr{
			Operator: scan.Token{scan.Minus, "-", nil, 1},
			Right:    &ast.LiteralExpr{123},
		},
		Operator: scan.Token{scan.Star, "*", nil, 1},
		Right:    &ast.GroupingExpr{&ast.LiteralExpr{45.67}},
	}
	printer := ast.Printer{}
	fmt.Println(printer.Print(binaryExpr))
}
