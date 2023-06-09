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
	ErrMissingIfLeftParen
	ErrMissingWhileLeftParen
	ErrMissingIfRightParen
	ErrMissingConditionRightParen
	ErrMissingForLeftParen
	ErrMissingForRightParen
	ErrMissingConditionSemicolon
	ErrArgumentsRightParen
	ErrFunctionName
	ErrFunctionLeftParen
	ErrFunctionRightParen
	ErrParameterName
	ErrFunctionLeftBrace
	ErrReturnSemicolon
	ErrClassName
	ErrClassLeftBrace
	ErrClassRightBrace
	ErrClassProperty
	ErrSuperclassName
	ErrSuperclassDot
	ErrSuperclassMethod
)

var errorMessages = map[ErrorKind]string{
	ErrMissingRightParen:          "missing ')' after expression",
	ErrMissingRightBrace:          "missing '}' after block",
	ErrMissingValueSemicolon:      "missing ';' after value",
	ErrMissingExprSemicolon:       "missing ';' after expression",
	ErrMissingVarSemicolon:        "missing ';' after variable declaration",
	ErrMissingExpr:                "missing expression",
	ErrMissingVariableName:        "missing variable name",
	ErrInvalidAssignTarget:        "invalid assignment target",
	ErrMissingIfLeftParen:         "missing '(' after 'if'",
	ErrMissingIfRightParen:        "missing ')' after 'if' condition",
	ErrMissingWhileLeftParen:      "missing '(' after 'while'",
	ErrMissingConditionRightParen: "missing ')' after condition",
	ErrMissingForLeftParen:        "missing '(' after 'for'",
	ErrMissingForRightParen:       "missing ')' after for clauses",
	ErrMissingConditionSemicolon:  "missing ';' after loop condition",
	ErrArgumentsRightParen:        "missing ')' after arguments",
	ErrFunctionName:               "missing function name",
	ErrFunctionLeftParen:          "missing '(' after function name",
	ErrFunctionRightParen:         "missing ')' after parameters",
	ErrParameterName:              "missing parameter name",
	ErrFunctionLeftBrace:          "missing '{' before function body",
	ErrReturnSemicolon:            "missing ';' after return value",
	ErrClassName:                  "missing class name",
	ErrClassLeftBrace:             "missing '{' before class body",
	ErrClassRightBrace:            "missing '}' after class body",
	ErrClassProperty:              "missing property name after '.'",
	ErrSuperclassName:             "missing superclass name",
	ErrSuperclassDot:              "missing '.' after 'super'",
	ErrSuperclassMethod:           "missing superclass method name",
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
