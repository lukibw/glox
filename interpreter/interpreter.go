package interpreter

import (
	"fmt"
	"strings"

	"github.com/lukibw/glox/ast"
)

type Interpreter struct {
	globals *env
	env     *env
	locals  map[ast.Expr]int
}

func New(locals map[ast.Expr]int) *Interpreter {
	globals := newEnv(nil)
	globals.define("clock", clock{})
	return &Interpreter{globals, globals, locals}
}

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	switch e := expr.(type) {
	case *ast.AssignExpr:
		return i.assignExpr(e)
	case *ast.BinaryExpr:
		return i.binaryExpr(e)
	case *ast.CallExpr:
		return i.callExpr(e)
	case *ast.GroupingExpr:
		return i.groupingExpr(e)
	case *ast.LogicalExpr:
		return i.logicalExpr(e)
	case *ast.LiteralExpr:
		return i.literalExpr(e)
	case *ast.UnaryExpr:
		return i.unaryExpr(e)
	case *ast.VarExpr:
		return i.varExpr(e)
	default:
		panic(fmt.Sprintf("interpreter: cannot evaluate an expression of type %T", e))
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) (any, error) {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		return i.blockStmt(s)
	case *ast.ExpressionStmt:
		return i.expressionStmt(s)
	case *ast.FunctionStmt:
		return i.functionStmt(s)
	case *ast.IfStmt:
		return i.ifStmt(s)
	case *ast.PrintStmt:
		return i.printStmt(s)
	case *ast.ReturnStmt:
		return i.returnStmt(s)
	case *ast.WhileStmt:
		return i.whileStmt(s)
	case *ast.VarStmt:
		return i.varStmt(s)
	default:
		panic(fmt.Sprintf("resolver: cannot execute a statement of type %T", s))
	}
}

func (i *Interpreter) executeBlock(stmts []ast.Stmt, env *env) (any, error) {
	previous := i.env
	i.env = env
	defer func() {
		i.env = previous
	}()
	for _, stmt := range stmts {
		value, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
		if value != nil {
			return value, nil
		}
	}
	return nil, nil
}

func (i *Interpreter) lookUpVariable(name ast.Token, expr ast.Expr) (any, error) {
	distance, ok := i.locals[expr]
	if !ok {
		return i.globals.get(name)
	}
	return i.env.getAt(distance, name)
}

func (i *Interpreter) Run(stmts []ast.Stmt) error {
	var err error
	for _, stmt := range stmts {
		if _, err = i.execute(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) assignExpr(expr *ast.AssignExpr) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	distance, ok := i.locals[expr]
	if !ok {
		i.globals.assign(expr.Name, value)
	} else {
		i.env.assignAt(distance, expr.Name, value)
	}
	return value, nil
}

func (i *Interpreter) binaryExpr(expr *ast.BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Kind {
	case ast.Plus:
		left_str, left_ok := left.(string)
		right_str, right_ok := right.(string)
		if left_ok && right_ok {
			sb := strings.Builder{}
			sb.WriteString(left_str)
			sb.WriteString(right_str)
			return sb.String(), nil
		}
		left_num, left_ok := left.(float64)
		right_num, right_ok := right.(float64)
		if left_ok && right_ok {
			return left_num + right_num, nil
		}
		return nil, &Error{expr.Operator, ErrNumberOrStringOperands}
	case ast.Minus:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case ast.Slash:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case ast.Star:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case ast.Greater:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case ast.GreaterEqual:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case ast.Less:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case ast.LessEqual:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case ast.BangEqual:
		return left != right, nil
	case ast.EqualEqual:
		return left == right, nil
	}
	panic("interpreter: cannot match operator for binary expression")
}

func (i *Interpreter) callExpr(expr *ast.CallExpr) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	arguments := make([]any, len(expr.Arguments))
	for k := 0; k < len(expr.Arguments); k++ {
		arguments[k], err = i.evaluate(expr.Arguments[k])
		if err != nil {
			return nil, err
		}
	}
	fn, ok := callee.(callable)
	if !ok {
		return nil, &Error{expr.Paren, ErrFunctionOrClassCallable}
	}
	if fn.arity() > len(arguments) {
		return nil, &Error{expr.Paren, ErrFunctionTooFewArgs}
	}
	if fn.arity() < len(arguments) {
		return nil, &Error{expr.Paren, ErrFunctionTooManyArgs}
	}
	return fn.call(i, arguments)
}

func (i *Interpreter) groupingExpr(expr *ast.GroupingExpr) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) logicalExpr(expr *ast.LogicalExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Kind == ast.Or {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) literalExpr(expr *ast.LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) unaryExpr(expr *ast.UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Kind {
	case ast.Minus:
		if err := checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case ast.Bang:
		return !isTruthy(right), nil
	}
	panic("interpreter: cannot match operator for unary expression")
}

func (i *Interpreter) varExpr(expr *ast.VarExpr) (any, error) {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) blockStmt(stmt *ast.BlockStmt) (any, error) {
	return i.executeBlock(stmt.Statements, newEnv(i.env))
}

func (i *Interpreter) expressionStmt(stmt *ast.ExpressionStmt) (any, error) {
	_, err := i.evaluate(stmt.Expression)
	return nil, err
}

func (i *Interpreter) functionStmt(stmt *ast.FunctionStmt) (any, error) {
	i.env.define(stmt.Name.Lexeme, &function{stmt, i.env})
	return nil, nil
}

func (i *Interpreter) ifStmt(stmt *ast.IfStmt) (any, error) {
	value, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(value) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) printStmt(stmt *ast.PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(stringify(value))
	return nil, nil
}

func (i *Interpreter) returnStmt(stmt *ast.ReturnStmt) (any, error) {
	if stmt.Value != nil {
		return i.evaluate(stmt.Value)
	}
	return nil, nil
}

func (i *Interpreter) whileStmt(stmt *ast.WhileStmt) (any, error) {
	for {
		value, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !isTruthy(value) {
			break
		}
		value, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
		if value != nil {
			return value, nil
		}
	}
	return nil, nil
}

func (i *Interpreter) varStmt(stmt *ast.VarStmt) (any, error) {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.env.define(stmt.Name.Lexeme, value)
	return nil, nil
}
