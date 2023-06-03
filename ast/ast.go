package ast

import "github.com/lukibw/glox/token"

type Expr[T any] interface {
	Accept(ExprVisitor[T]) (T, error)
}

type ExprVisitor[T any] interface {
	VisitBinaryExpr(*BinaryExpr[T]) (T, error)
	VisitGroupingExpr(*GroupingExpr[T]) (T, error)
	VisitLiteralExpr(*LiteralExpr[T]) (T, error)
	VisitUnaryExpr(*UnaryExpr[T]) (T, error)
}

type BinaryExpr[T any] struct {
	Left     Expr[T]
	Operator token.Token
	Right    Expr[T]
}

func (e *BinaryExpr[T]) Accept(v ExprVisitor[T]) (T, error) {
	return v.VisitBinaryExpr(e)
}

type GroupingExpr[T any] struct {
	Expression Expr[T]
}

func (e *GroupingExpr[T]) Accept(v ExprVisitor[T]) (T, error) {
	return v.VisitGroupingExpr(e)
}

type LiteralExpr[T any] struct {
	Value any
}

func (e *LiteralExpr[T]) Accept(v ExprVisitor[T]) (T, error) {
	return v.VisitLiteralExpr(e)
}

type UnaryExpr[T any] struct {
	Operator token.Token
	Right    Expr[T]
}

func (e *UnaryExpr[T]) Accept(v ExprVisitor[T]) (T, error) {
	return v.VisitUnaryExpr(e)
}
