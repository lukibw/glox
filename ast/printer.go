package ast

import (
	"fmt"
	"strings"
)

type Printer struct {
	builder strings.Builder
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(expr Expr) string {
	p.builder.Reset()
	expr.Accept(p)
	return p.builder.String()
}

func (p *Printer) parenthesize(name string, exprs ...Expr) {
	p.builder.WriteRune('(')
	p.builder.WriteString(name)
	for _, expr := range exprs {
		p.builder.WriteRune(' ')
		expr.Accept(p)
	}
	p.builder.WriteRune(')')
}

func (p *Printer) VisitBinaryExpr(expr *BinaryExpr) {
	p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *Printer) VisitGroupingExpr(expr *GroupingExpr) {
	p.parenthesize("group", expr.Expression)
}

func (p *Printer) VisitLiteralExpr(expr *LiteralExpr) {
	if expr.Value == nil {
		p.builder.WriteString("nil")
	} else {
		p.builder.WriteString(fmt.Sprint(expr.Value))
	}
}

func (p *Printer) VisitUnaryExpr(expr *UnaryExpr) {
	p.parenthesize(expr.Operator.Lexeme, expr.Right)
}
