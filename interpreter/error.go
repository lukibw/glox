package interpreter

import (
	"fmt"

	"github.com/lukibw/glox/ast"
)

type ErrorKind int

const (
	ErrNumberOperand ErrorKind = iota
	ErrNumberOperands
	ErrNumberOrStringOperands
	ErrUndefinedVariable
)

var runtimeErrorMessages = map[ErrorKind]string{
	ErrNumberOperand:          "operand must be a number",
	ErrNumberOperands:         "operands must be numbers",
	ErrNumberOrStringOperands: "operands must be two numbers or two strings",
	ErrUndefinedVariable:      "undefined variable",
}

func (k ErrorKind) String() string {
	return runtimeErrorMessages[k]
}

type Error struct {
	Token ast.Token
	Kind  ErrorKind
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.Kind, e.Token.Line)
}
