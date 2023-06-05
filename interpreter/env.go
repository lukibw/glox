package interpreter

import (
	"github.com/lukibw/glox/ast"
)

type env struct {
	values    map[string]any
	enclosing *env
}

func newEnv(enclosing *env) *env {
	return &env{make(map[string]any), enclosing}
}

func (e *env) define(name string, value any) {
	e.values[name] = value
}

func (e *env) get(name ast.Token) (any, error) {
	v, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		return nil, &Error{name, ErrUndefinedVariable}
	}
	return v, nil
}

func (e *env) assign(name ast.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; !ok {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}
		return &Error{name, ErrUndefinedVariable}
	}
	e.values[name.Lexeme] = value
	return nil
}
