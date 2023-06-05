package ast

import (
	"fmt"

	"github.com/lukibw/glox/token"
)

type RuntimeErrorKind int

const (
	ErrNumberOperand RuntimeErrorKind = iota
	ErrNumberOperands
	ErrNumberOrStringOperands
	ErrUndefinedVariable
)

var runtimeErrorMessages = map[RuntimeErrorKind]string{
	ErrNumberOperand:          "operand must be a number",
	ErrNumberOperands:         "operands must be numbers",
	ErrNumberOrStringOperands: "operands must be two numbers or two strings",
	ErrUndefinedVariable:      "undefined variable",
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
