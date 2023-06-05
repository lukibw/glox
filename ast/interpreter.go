package ast

import (
	"fmt"
	"strings"

	"github.com/lukibw/glox/token"
)

type Interpreter struct {
	env *Env
}

func NewInterpreter() *Interpreter {
	return &Interpreter{NewEnv(nil)}
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

func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &RuntimeError{operator, ErrNumberOperand}
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left, right any) error {
	_, left_ok := left.(float64)
	_, right_ok := right.(float64)
	if left_ok && right_ok {
		return nil
	}
	return &RuntimeError{operator, ErrNumberOperands}
}

func (i *Interpreter) evaluate(expr Expr[any]) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt[any]) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt[any], env *Env) (any, error) {
	var err error
	previous := i.env
	i.env = env
	for _, statement := range statements {
		_, err = i.execute(statement)
		if err != nil {
			break
		}
	}
	i.env = previous
	return nil, err
}

func (i *Interpreter) Interpret(statements []Stmt[any]) error {
	for _, statement := range statements {
		if _, err := i.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr[any]) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	i.env.Assign(expr.Name, value)
	return value, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr[any]) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Kind {
	case token.Plus:
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
		return nil, &RuntimeError{expr.Operator, ErrNumberOrStringOperands}
	case token.Minus:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.Slash:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.Star:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.Greater:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GreaterEqual:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.Less:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LessEqual:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.BangEqual:
		return left != right, nil
	case token.EqualEqual:
		return left == right, nil
	}
	panic("fix this")
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr[any]) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr[any]) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr[any]) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Kind {
	case token.Minus:
		if err := i.checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case token.Bang:
		return !i.isTruthy(right), nil
	}
	panic("fix this")
}

func (i *Interpreter) VisitVariableExpr(expr *VariableExpr[any]) (any, error) {
	return i.env.Get(expr.Name)
}

func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt[any]) (any, error) {
	_, err := i.evaluate(stmt.Expression)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt[any]) (any, error) {
	v, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(i.stringify(v))
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt[any]) (any, error) {
	return i.executeBlock(stmt.Statements, NewEnv(i.env))
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt[any]) (any, error) {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.env.Define(stmt.Name.Lexeme, value)
	return nil, nil
}
