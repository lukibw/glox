package parser

import (
	"fmt"

	"github.com/lukibw/glox/token"
)

type ErrorKind int

const (
	ErrMissingRightParen ErrorKind = iota
	ErrMissingExpr
)

var errorMessages = map[ErrorKind]string{
	ErrMissingRightParen: "expected ')' after expression",
	ErrMissingExpr:       "expected expression",
}

func (k ErrorKind) String() string {
	return errorMessages[k]
}

type Error struct {
	Token token.Token
	Kind  ErrorKind
}

func (e *Error) Error() string {
	if e.Token.Kind == token.Eof {
		return fmt.Sprintf("[line %d] parsing error at end: %s", e.Token.Line, e.Kind)
	}
	return fmt.Sprintf("[line %d] parsing error at '%s': %s", e.Token.Line, e.Token.Lexeme, e.Kind)
}
