package ast

import "github.com/lukibw/glox/token"

type Expr[T any] interface {
	Accept(ExprVisitor[T]) (T, error)
}

type ExprVisitor[T any] interface {
	VisitAssignExpr(*AssignExpr[T]) (T, error)
	VisitBinaryExpr(*BinaryExpr[T]) (T, error)
	VisitGroupingExpr(*GroupingExpr[T]) (T, error)
	VisitLiteralExpr(*LiteralExpr[T]) (T, error)
	VisitUnaryExpr(*UnaryExpr[T]) (T, error)
	VisitVariableExpr(*VariableExpr[T]) (T, error)
}

type AssignExpr[T any] struct {
	Name  token.Token
	Value Expr[T]
}

func (e *AssignExpr[T]) Accept(v ExprVisitor[T]) (T, error) {
	return v.VisitAssignExpr(e)
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

type VariableExpr[T any] struct {
	Name token.Token
}

func (e *VariableExpr[T]) Accept(v ExprVisitor[T]) (T, error) {
	return v.VisitVariableExpr(e)
}

type Stmt[T any] interface {
	Accept(StmtVisitor[T]) (T, error)
}

type StmtVisitor[T any] interface {
	VisitBlockStmt(*BlockStmt[T]) (T, error)
	VisitExpressionStmt(*ExpressionStmt[T]) (T, error)
	VisitPrintStmt(*PrintStmt[T]) (T, error)
	VisitVarStmt(*VarStmt[T]) (T, error)
}

type BlockStmt[T any] struct {
	Statements []Stmt[T]
}

func (s *BlockStmt[T]) Accept(v StmtVisitor[T]) (T, error) {
	return v.VisitBlockStmt(s)
}

type ExpressionStmt[T any] struct {
	Expression Expr[T]
}

func (s *ExpressionStmt[T]) Accept(v StmtVisitor[T]) (T, error) {
	return v.VisitExpressionStmt(s)
}

type PrintStmt[T any] struct {
	Expression Expr[T]
}

func (s *PrintStmt[T]) Accept(v StmtVisitor[T]) (T, error) {
	return v.VisitPrintStmt(s)
}

type VarStmt[T any] struct {
	Name        token.Token
	Initializer Expr[T]
}

func (s *VarStmt[T]) Accept(v StmtVisitor[T]) (T, error) {
	return v.VisitVarStmt(s)
}
