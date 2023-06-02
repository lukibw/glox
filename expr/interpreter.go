package expr

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

type Interpreter struct {
	counter int
	value   any
	issue   error
}

func NewInterpreter() *Interpreter {
	return &Interpreter{0, nil, nil}
}

func (i *Interpreter) newIssue(operator token.Token, kind RuntimeErrorKind) {
	i.issue = &RuntimeError{operator, kind}
	i.value = nil
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

func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) bool {
	if _, ok := operand.(float64); ok {
		return true
	}
	i.newIssue(operator, ErrNumberOperand)
	return false
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left, right any) bool {
	_, left_ok := left.(float64)
	_, right_ok := right.(float64)
	if left_ok && right_ok {
		return true
	}
	i.newIssue(operator, ErrNumberOperands)
	return false
}

func (i *Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprint(value)
}

func (i *Interpreter) Interpret(exp Expr) (string, error) {
	value := i.evaluate(exp)
	if i.issue != nil {
		return "", i.issue
	} else {
		return i.stringify(value), nil
	}
}

func (i *Interpreter) evaluate(exp Expr) any {
	exp.Accept(i)
	return i.value
}

func (i *Interpreter) VisitBinary(exp *Binary) {
	if i.issue != nil {
		return
	}
	left := i.evaluate(exp.Left)
	right := i.evaluate(exp.Right)
	switch exp.Operator.Kind {
	case token.Plus:
		left_str, left_ok := left.(string)
		right_str, right_ok := right.(string)
		if left_ok && right_ok {
			sb := strings.Builder{}
			sb.WriteString(left_str)
			sb.WriteString(right_str)
			i.value = sb.String()
			return
		}
		left_num, left_ok := left.(float64)
		right_num, right_ok := right.(float64)
		if left_ok && right_ok {
			i.value = left_num + right_num
			return
		}
		i.newIssue(exp.Operator, ErrNumberOrStringOperands)
	case token.Minus:
		if i.checkNumberOperands(exp.Operator, left, right) {
			i.value = left.(float64) - right.(float64)
		}
	case token.Slash:
		if i.checkNumberOperands(exp.Operator, left, right) {
			i.value = left.(float64) / right.(float64)
		}
	case token.Star:
		if i.checkNumberOperands(exp.Operator, left, right) {
			i.value = left.(float64) * right.(float64)
		}
	case token.Greater:
		if i.checkNumberOperands(exp.Operator, left, right) {
			i.value = left.(float64) > right.(float64)
		}
	case token.GreaterEqual:
		if i.checkNumberOperands(exp.Operator, left, right) {
			i.value = left.(float64) >= right.(float64)
		}
	case token.Less:
		if i.checkNumberOperands(exp.Operator, left, right) {
			i.value = left.(float64) < right.(float64)
		}
	case token.LessEqual:
		if i.checkNumberOperands(exp.Operator, left, right) {
			i.value = left.(float64) <= right.(float64)
		}
	case token.BangEqual:
		i.value = left != right
	case token.EqualEqual:
		i.value = left == right
	}
}

func (i *Interpreter) VisitGrouping(exp *Grouping) {
	if i.issue != nil {
		return
	}
	i.evaluate(exp.Expression)
}

func (i *Interpreter) VisitLiteral(exp *Literal) {
	if i.issue != nil {
		return
	}
	i.value = exp.Value
}

func (i *Interpreter) VisitUnary(exp *Unary) {
	if i.issue != nil {
		return
	}
	right := i.evaluate(exp.Right)
	switch exp.Operator.Kind {
	case token.Minus:
		i.checkNumberOperand(exp.Operator, right)
		i.value = -right.(float64)
	case token.Bang:
		i.value = !i.isTruthy(right)
	}
}
