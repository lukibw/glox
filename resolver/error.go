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
	ErrSelfInheritClass
	ErrSuperOutsideClass
	ErrSuperNoSuperclass
)

var errorMessages = map[ErrorKind]string{
	ErrVarInitializer:    "cannot read local variable in its own initializer",
	ErrVarDuplicate:      "cannot declare a variable that is already in this scope",
	ErrTopLevelReturn:    "cannot return from top-level code",
	ErrInitializerReturn: "cannot return a value from an initializer",
	ErrThisOutsideClass:  "cannot use 'this' outside of a class",
	ErrSelfInheritClass:  "a class cannot inherit from itself",
	ErrSuperOutsideClass: "cannot use 'super' outside of a class",
	ErrSuperNoSuperclass: "cannot use 'super' in a class with no superclass",
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
