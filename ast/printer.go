package ast

import (
	"fmt"
	"strings"
)

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) parenthesize(name string, exprs ...Expr[string]) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteRune(' ')
		s, err := expr.Accept(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
	}
	sb.WriteRune(')')
	return sb.String(), nil
}

func (p *Printer) String(expr Expr[string]) (string, error) {
	return expr.Accept(p)
}

func (p *Printer) VisitBinaryExpr(expr *BinaryExpr[string]) (string, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *Printer) VisitGroupingExpr(expr *GroupingExpr[string]) (string, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *Printer) VisitLiteralExpr(expr *LiteralExpr[string]) (string, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprint(expr.Value), nil
}

func (p *Printer) VisitUnaryExpr(expr *UnaryExpr[string]) (string, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}
