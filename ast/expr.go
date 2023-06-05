package ast

type Expr interface {
	Accept(ExprVisitor) (any, error)
}

type ExprVisitor interface {
	VisitAssignExpr(*AssignExpr) (any, error)
	VisitBinaryExpr(*BinaryExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitVarExpr(*VarExpr) (any, error)
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (e *AssignExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitAssignExpr(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitUnaryExpr(e)
}

type VarExpr struct {
	Name Token
}

func (e *VarExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitVarExpr(e)
}
