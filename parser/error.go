package parser

import (
	"fmt"

	"github.com/lukibw/glox/ast"
)

type ErrorKind int

const (
	ErrMissingRightParen ErrorKind = iota
	ErrMissingRightBrace
	ErrMissingValueSemicolon
	ErrMissingExprSemicolon
	ErrMissingVarSemicolon
	ErrMissingExpr
	ErrMissingVariableName
	ErrInvalidAssignTarget
)

var errorMessages = map[ErrorKind]string{
	ErrMissingRightParen:     "missing ')' after expression",
	ErrMissingRightBrace:     "missing '}' after block",
	ErrMissingValueSemicolon: "missing ';' after value",
	ErrMissingExprSemicolon:  "missing ';' after expression",
	ErrMissingVarSemicolon:   "missing ';' after variable declaration",
	ErrMissingExpr:           "missing expression",
	ErrMissingVariableName:   "missing variable name",
	ErrInvalidAssignTarget:   "invalid assignment target",
}

func (k ErrorKind) String() string {
	return errorMessages[k]
}

type Error struct {
	Token ast.Token
	Kind  ErrorKind
}

func (e *Error) Error() string {
	var where string
	if e.Token.Kind == ast.Eof {
		where = "end"
	} else {
		where = fmt.Sprintf("'%s'", e.Token.Lexeme)
	}
	return fmt.Sprintf("[line %d] error at %s: %s", e.Token.Line, where, e.Kind)
}
