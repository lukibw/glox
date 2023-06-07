package ast

type Expr interface {
	expr()
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (e *AssignExpr) expr() {}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *BinaryExpr) expr() {}

type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func (e *CallExpr) expr() {}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) expr() {}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) expr() {}

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *LogicalExpr) expr() {}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (e *UnaryExpr) expr() {}

type VarExpr struct {
	Name Token
}

func (e *VarExpr) expr() {}
