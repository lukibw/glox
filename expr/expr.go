package expr

import "github.com/lukibw/glox/token"

type Expr interface {
	Accept(Visitor)
}

type Visitor interface {
	VisitBinary(*Binary)
	VisitGrouping(*Grouping)
	VisitLiteral(*Literal)
	VisitUnary(*Unary)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *Binary) Accept(v Visitor) {
	v.VisitBinary(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(v Visitor) {
	v.VisitGrouping(e)
}

type Literal struct {
	Value any
}

func (e *Literal) Accept(v Visitor) {
	v.VisitLiteral(e)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (e *Unary) Accept(v Visitor) {
	v.VisitUnary(e)
}
