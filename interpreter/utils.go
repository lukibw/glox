package interpreter

import (
	"fmt"

	"github.com/lukibw/glox/ast"
)

func stringify(value any) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprint(value)
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if v, ok := value.(bool); ok {
		return v
	}
	return true
}

func checkNumberOperand(operator ast.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &Error{operator, ErrNumberOperand}
}

func checkNumberOperands(operator ast.Token, left, right any) error {
	_, left_ok := left.(float64)
	_, right_ok := right.(float64)
	if left_ok && right_ok {
		return nil
	}
	return &Error{operator, ErrNumberOperands}
}
