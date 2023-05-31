package ast

import (
	"github.com/lukibw/glox/scan"
)

type Expr interface {
	Accept(ExprVisitor)
}

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr)
	VisitGroupingExpr(*GroupingExpr)
	VisitLiteralExpr(*LiteralExpr)
	VisitUnaryExpr(*UnaryExpr)
}

type BinaryExpr struct {
	Left     Expr
	Operator scan.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) {
	v.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) {
	v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(v ExprVisitor) {
	v.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator scan.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) {
	v.VisitUnaryExpr(e)
}
