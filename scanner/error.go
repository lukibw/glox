package scanner

import "fmt"

type ErrorKind int

const (
	ErrUnterminatedString ErrorKind = iota
	ErrUnexpectedCharacter
)

var errorMessages = map[ErrorKind]string{
	ErrUnterminatedString:  "unterminated string",
	ErrUnexpectedCharacter: "unexpected character",
}

func (k ErrorKind) String() string {
	return errorMessages[k]
}

type Error struct {
	Line int
	Kind ErrorKind
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] error: %s", e.Line, e.Kind)
}
