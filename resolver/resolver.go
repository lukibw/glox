package resolver

import (
	"fmt"

	"github.com/lukibw/glox/ast"
)

type functionKind int

const (
	noFunction functionKind = iota
	normalFunction
	methodFunction
	initializerFunction
)

type classKind int

const (
	noClass classKind = iota
	normalClass
	subClass
)

type Resolver struct {
	scopes       *stack
	locals       map[ast.Expr]int
	currentFn    functionKind
	currentClass classKind
}

func New() *Resolver {
	return &Resolver{newStack(), make(map[ast.Expr]int), noFunction, noClass}
}

func (r *Resolver) declare(name ast.Token) error {
	if r.scopes.isEmpty() {
		return nil
	}
	scope := r.scopes.peek()
	if _, ok := scope[name.Lexeme]; ok {
		return &Error{name, ErrVarDuplicate}
	}
	r.scopes.peek()[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name ast.Token) {
	if !r.scopes.isEmpty() {
		r.scopes.peek()[name.Lexeme] = true
	}
}

func (r *Resolver) beginScope() {
	r.scopes.push()
}

func (r *Resolver) endScope() {
	r.scopes.pop()
}

func (r *Resolver) resolveExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.AssignExpr:
		return r.assignExpr(e)
	case *ast.BinaryExpr:
		return r.binaryExpr(e)
	case *ast.CallExpr:
		return r.callExpr(e)
	case *ast.GetExpr:
		return r.getExpr(e)
	case *ast.GroupingExpr:
		return r.groupingExpr(e)
	case *ast.LiteralExpr:
		return r.literalExpr(e)
	case *ast.LogicalExpr:
		return r.logicalExpr(e)
	case *ast.SetExpr:
		return r.setExpr(e)
	case *ast.SuperExpr:
		return r.superExpr(e)
	case *ast.ThisExpr:
		return r.thisExpr(e)
	case *ast.UnaryExpr:
		return r.unaryExpr(e)
	case *ast.VarExpr:
		return r.varExpr(e)
	default:
		panic(fmt.Sprintf("resolver: cannot resolve an expression of type %T", e))
	}
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		return r.blockStmt(s)
	case *ast.ClassStmt:
		return r.classStmt(s)
	case *ast.ExpressionStmt:
		return r.expressionStmt(s)
	case *ast.FunctionStmt:
		return r.functionStmt(s)
	case *ast.IfStmt:
		return r.ifStmt(s)
	case *ast.PrintStmt:
		return r.printStmt(s)
	case *ast.ReturnStmt:
		return r.returnStmt(s)
	case *ast.WhileStmt:
		return r.whileStmt(s)
	case *ast.VarStmt:
		return r.varStmt(s)
	default:
		panic(fmt.Sprintf("resolver: cannot resolve a statement of type %T", s))
	}
}

func (r *Resolver) resolveStmts(stmts []ast.Stmt) error {
	var err error
	for _, stmt := range stmts {
		if err = r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveLocal(expr ast.Expr, name ast.Token) error {
	for i := r.scopes.size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.get(i)[name.Lexeme]; ok {
			r.locals[expr] = r.scopes.size() - 1 - i
		}
	}
	return nil
}

func (r *Resolver) resolveFunction(fn *ast.FunctionStmt, kind functionKind) error {
	enclosingFn := r.currentFn
	r.currentFn = kind
	r.beginScope()
	defer func() {
		r.endScope()
		r.currentFn = enclosingFn
	}()
	for _, param := range fn.Params {
		if err := r.declare(param); err != nil {
			return err
		}
		r.define(param)
	}
	return r.resolveStmts(fn.Body)
}

func (r *Resolver) Run(stmts []ast.Stmt) (map[ast.Expr]int, error) {
	if err := r.resolveStmts(stmts); err != nil {
		return nil, err
	}
	return r.locals, nil
}

func (r *Resolver) assignExpr(expr *ast.AssignExpr) error {
	if err := r.resolveExpr(expr.Value); err != nil {
		return err
	}
	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) binaryExpr(expr *ast.BinaryExpr) error {
	if err := r.resolveExpr(expr.Left); err != nil {
		return err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) callExpr(expr *ast.CallExpr) error {
	err := r.resolveExpr(expr.Callee)
	if err != nil {
		return err
	}
	for _, arg := range expr.Arguments {
		if err = r.resolveExpr(arg); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) getExpr(expr *ast.GetExpr) error {
	return r.resolveExpr(expr.Object)
}

func (r *Resolver) groupingExpr(expr *ast.GroupingExpr) error {
	return r.resolveExpr(expr.Expression)
}

func (r *Resolver) literalExpr(expr *ast.LiteralExpr) error {
	return nil
}

func (r *Resolver) logicalExpr(expr *ast.LogicalExpr) error {
	if err := r.resolveExpr(expr.Left); err != nil {
		return err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) setExpr(expr *ast.SetExpr) error {
	if err := r.resolveExpr(expr.Value); err != nil {
		return err
	}
	return r.resolveExpr(expr.Object)
}

func (r *Resolver) superExpr(expr *ast.SuperExpr) error {
	if r.currentClass == noClass {
		return &Error{expr.Keyword, ErrSuperOutsideClass}
	}
	if r.currentClass != subClass {
		return &Error{expr.Keyword, ErrSuperNoSuperclass}
	}
	return r.resolveLocal(expr, expr.Keyword)
}

func (r *Resolver) thisExpr(expr *ast.ThisExpr) error {
	if r.currentClass == noClass {
		return &Error{expr.Keyword, ErrThisOutsideClass}
	}
	return r.resolveLocal(expr, expr.Keyword)
}

func (r *Resolver) unaryExpr(expr *ast.UnaryExpr) error {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) varExpr(expr *ast.VarExpr) error {
	if !r.scopes.isEmpty() {
		if defined, ok := r.scopes.peek()[expr.Name.Lexeme]; ok && !defined {
			return &Error{expr.Name, ErrVarInitializer}
		}
	}
	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) blockStmt(stmt *ast.BlockStmt) error {
	r.beginScope()
	defer r.endScope()
	return r.resolveStmts(stmt.Statements)
}

func (r *Resolver) classStmt(stmt *ast.ClassStmt) error {
	enclosingClass := r.currentClass
	r.currentClass = normalClass
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}
	r.define(stmt.Name)
	if stmt.Superclass != nil {
		r.currentClass = subClass
		if stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
			return &Error{stmt.Superclass.Name, ErrSelfInheritClass}
		}
		if err = r.resolveExpr(stmt.Superclass); err != nil {
			return err
		}
		r.beginScope()
		r.scopes.peek()["super"] = true
	}
	r.beginScope()
	r.scopes.peek()["this"] = true
	for _, method := range stmt.Methods {
		kind := methodFunction
		if method.Name.Lexeme == "init" {
			kind = initializerFunction
		}
		if err = r.resolveFunction(method, kind); err != nil {
			return err
		}
	}
	r.endScope()
	if stmt.Superclass != nil {
		r.endScope()
	}
	r.currentClass = enclosingClass
	return nil
}

func (r *Resolver) expressionStmt(stmt *ast.ExpressionStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) functionStmt(stmt *ast.FunctionStmt) error {
	if err := r.declare(stmt.Name); err != nil {
		return err
	}
	r.define(stmt.Name)
	return r.resolveFunction(stmt, normalFunction)
}

func (r *Resolver) ifStmt(stmt *ast.IfStmt) error {
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	if err = r.resolveStmt(stmt.ThenBranch); err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		if err = r.resolveStmt(stmt.ElseBranch); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) printStmt(stmt *ast.PrintStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) returnStmt(stmt *ast.ReturnStmt) error {
	if r.currentFn == noFunction {
		return &Error{stmt.Keyword, ErrTopLevelReturn}
	}
	if stmt.Value != nil {
		if r.currentFn == initializerFunction {
			return &Error{stmt.Keyword, ErrInitializerReturn}
		}
		return r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) whileStmt(stmt *ast.WhileStmt) error {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return err
	}
	return r.resolveStmt(stmt.Body)
}

func (r *Resolver) varStmt(stmt *ast.VarStmt) error {
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}
	if stmt.Initializer != nil {
		if err = r.resolveExpr(stmt.Initializer); err != nil {
			return err
		}
	}
	r.define(stmt.Name)
	return nil
}
