package ast

type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	VisitBlockStmt(*BlockStmt) error
	VisitExpressionStmt(*ExpressionStmt) error
	VisitPrintStmt(*PrintStmt) error
	VisitVarStmt(*VarStmt) error
	VisitIfStmt(*IfStmt) error
	VisitWhileStmt(*WhileStmt) error
}

type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) Accept(v StmtVisitor) error {
	return v.VisitBlockStmt(s)
}

type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) error {
	return v.VisitExpressionStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) Accept(v StmtVisitor) error {
	return v.VisitPrintStmt(s)
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (s *VarStmt) Accept(v StmtVisitor) error {
	return v.VisitVarStmt(s)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (s *IfStmt) Accept(v StmtVisitor) error {
	return v.VisitIfStmt(s)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) Accept(v StmtVisitor) error {
	return v.VisitWhileStmt(s)
}
