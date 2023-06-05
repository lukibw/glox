package ast

import "github.com/lukibw/glox/token"

type Env struct {
	values    map[string]any
	enclosing *Env
}

func NewEnv(enclosing *Env) *Env {
	return &Env{make(map[string]any), enclosing}
}

func (e *Env) Define(name string, value any) {
	e.values[name] = value
}

func (e *Env) Get(name token.Token) (any, error) {
	v, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		return nil, &RuntimeError{name, ErrUndefinedVariable}
	}
	return v, nil
}

func (e *Env) Assign(name token.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(name, value)
		}
		return &RuntimeError{name, ErrUndefinedVariable}
	}
	e.values[name.Lexeme] = value
	return nil
}
