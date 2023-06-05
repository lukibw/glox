package interpreter

import (
	"fmt"
	"strings"

	"github.com/lukibw/glox/ast"
)

type Interpreter struct {
	globals *env
	env     *env
}

func New() *Interpreter {
	globals := newEnv(nil)
	globals.define("clock", Clock{})
	return &Interpreter{globals, globals}
}

func (i *Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprint(value)
}

func (i *Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if v, ok := value.(bool); ok {
		return v
	}
	return true
}

func (i *Interpreter) checkNumberOperand(operator ast.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &Error{operator, ErrNumberOperand}
}

func (i *Interpreter) checkNumberOperands(operator ast.Token, left, right any) error {
	_, left_ok := left.(float64)
	_, right_ok := right.(float64)
	if left_ok && right_ok {
		return nil
	}
	return &Error{operator, ErrNumberOperands}
}

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) error {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, env *env) error {
	var err error
	previous := i.env
	i.env = env
	for _, statement := range statements {
		err = i.execute(statement)
		if err != nil {
			break
		}
	}
	i.env = previous
	return err
}

func (i *Interpreter) Run(statements []ast.Stmt) error {
	var err error
	for _, statement := range statements {
		if err = i.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	i.env.assign(expr.Name, value)
	return value, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
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
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case ast.Slash:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case ast.Star:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case ast.Greater:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case ast.GreaterEqual:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case ast.Less:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case ast.LessEqual:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
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

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Kind {
	case ast.Minus:
		if err := i.checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case ast.Bang:
		return !i.isTruthy(right), nil
	}
	panic("interpreter: cannot match operator for unary expression")
}

func (i *Interpreter) VisitVarExpr(expr *ast.VarExpr) (any, error) {
	return i.env.get(expr.Name)
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Kind == ast.Or {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	arguments := make([]any, len(expr.Arguments))
	for k := 0; k < len(expr.Arguments); k++ {
		arg, err := i.evaluate(expr.Arguments[k])
		if err != nil {
			return nil, err
		}
		arguments[k] = arg
	}
	function, ok := callee.(Callable)
	if !ok {
		return nil, &Error{expr.Paren, ErrFunctionOrClassCallable}
	}
	if function.Arity() > len(arguments) {
		return nil, &Error{expr.Paren, ErrFunctionTooFewArgs}
	}
	if function.Arity() < len(arguments) {
		return nil, &Error{expr.Paren, ErrFunctionTooManyArgs}
	}
	value, err := function.Call(i, arguments)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) error {
	_, err := i.evaluate(stmt.Expression)
	return err
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) error {
	v, err := i.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Println(i.stringify(v))
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) error {
	return i.executeBlock(stmt.Statements, newEnv(i.env))
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) error {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	i.env.define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) error {
	value, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}
	if i.isTruthy(value) {
		if err = i.execute(stmt.ThenBranch); err != nil {
			return err
		}
	} else if stmt.ElseBranch != nil {
		if err = i.execute(stmt.ElseBranch); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) error {
	for {
		value, err := i.evaluate(stmt.Condition)
		if err != nil {
			return err
		}
		if !i.isTruthy(value) {
			break
		}
		if err = i.execute(stmt.Body); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) error {
	function := &Function{Declaration: stmt, Closure: i.env}
	i.env.define(stmt.Name.Lexeme, function)
	return nil
}

type returnError struct {
	value any
}

func (e *returnError) Error() string {
	return fmt.Sprint(e.value)
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) error {
	var value any
	var err error
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return err
		}
	}
	return &returnError{value}
}
