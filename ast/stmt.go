package ast

type Stmt interface {
	stmt()
}

type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) stmt() {}

type ClassStmt struct {
	Name    Token
	Methods []*FunctionStmt
}

func (s *ClassStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) stmt() {}

type FunctionStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (s *FunctionStmt) stmt() {}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (s *IfStmt) stmt() {}

type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) stmt() {}

type ReturnStmt struct {
	Keyword Token
	Value   Expr
}

func (s *ReturnStmt) stmt() {}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) stmt() {}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (s *VarStmt) stmt() {}
