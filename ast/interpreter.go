package ast

import (
	"fmt"
	"strings"

	"github.com/lukibw/glox/token"
)

type RuntimeErrorKind int

const (
	ErrNumberOperand RuntimeErrorKind = iota
	ErrNumberOperands
	ErrNumberOrStringOperands
)

var runtimeErrorMessages = map[RuntimeErrorKind]string{
	ErrNumberOperand:          "operand must be a number",
	ErrNumberOperands:         "operands must be numbers",
	ErrNumberOrStringOperands: "operands must be two numbers or two strings",
}

func (k RuntimeErrorKind) String() string {
	return runtimeErrorMessages[k]
}

type RuntimeError struct {
	Token token.Token
	Kind  RuntimeErrorKind
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.Kind, e.Token.Line)
}

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
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

func (i *Interpreter) Interpret(expr Expr[any]) (string, error) {
	v, err := i.evaluate(expr)
	if err != nil {
		return "", err
	}
	return i.stringify(v), nil
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
