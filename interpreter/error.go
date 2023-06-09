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
	ErrFunctionOrClassCallable
	ErrFunctionTooFewArgs
	ErrFunctionTooManyArgs
	ErrInstanceProperty
	ErrUndefinedProperty
	ErrSuperclassNotAClass
)

var runtimeErrorMessages = map[ErrorKind]string{
	ErrNumberOperand:           "operand must be a number",
	ErrNumberOperands:          "operands must be numbers",
	ErrNumberOrStringOperands:  "operands must be two numbers or two strings",
	ErrUndefinedVariable:       "undefined variable",
	ErrFunctionOrClassCallable: "callable must be a function or a class",
	ErrFunctionTooFewArgs:      "too few arguments passed to the function",
	ErrFunctionTooManyArgs:     "too many arguments passed to the function",
	ErrInstanceProperty:        "only instances have properties",
	ErrUndefinedProperty:       "undefined property",
	ErrSuperclassNotAClass:     "superclass must be a class",
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
