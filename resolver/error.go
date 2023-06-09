package resolver

import (
	"fmt"

	"github.com/lukibw/glox/ast"
)

type ErrorKind int

const (
	ErrVarInitializer ErrorKind = iota
	ErrVarDuplicate
	ErrTopLevelReturn
	ErrInitializerReturn
	ErrThisOutsideClass
)

var errorMessages = map[ErrorKind]string{
	ErrVarInitializer:    "cannot read local variable in its own initializer",
	ErrVarDuplicate:      "cannot declare a variable that is already in this scope",
	ErrTopLevelReturn:    "cannot return from top-level code",
	ErrInitializerReturn: "cannot return a value from an initializer",
	ErrThisOutsideClass:  "cannot use 'this' outside of a class",
}

func (k ErrorKind) String() string {
	return errorMessages[k]
}

type Error struct {
	Token ast.Token
	Kind  ErrorKind
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] error: %s", e.Token.Line, e.Kind)
}
