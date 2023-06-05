package interpreter

import (
	"fmt"
	"strings"

	"github.com/lukibw/glox/ast"
)

type Interpreter struct {
	env *env
}

func New() *Interpreter {
	return &Interpreter{newEnv(nil)}
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
