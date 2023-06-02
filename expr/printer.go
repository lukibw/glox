package expr

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

func (p *Printer) Print(exp Expr) string {
	p.builder.Reset()
	exp.Accept(p)
	return p.builder.String()
}

func (p *Printer) parenthesize(name string, exps ...Expr) {
	p.builder.WriteRune('(')
	p.builder.WriteString(name)
	for _, exp := range exps {
		p.builder.WriteRune(' ')
		exp.Accept(p)
	}
	p.builder.WriteRune(')')
}

func (p *Printer) VisitBinary(exp *Binary) {
	p.parenthesize(exp.Operator.Lexeme, exp.Left, exp.Right)
}

func (p *Printer) VisitGrouping(exp *Grouping) {
	p.parenthesize("group", exp.Expression)
}

func (p *Printer) VisitLiteral(exp *Literal) {
	if exp.Value == nil {
		p.builder.WriteString("nil")
	} else {
		p.builder.WriteString(fmt.Sprint(exp.Value))
	}
}

func (p *Printer) VisitUnary(exp *Unary) {
	p.parenthesize(exp.Operator.Lexeme, exp.Right)
}
