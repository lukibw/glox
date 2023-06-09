package interpreter

import (
	"fmt"

	"github.com/lukibw/glox/ast"
)

type class struct {
	name    string
	methods map[string]*function
}

func (c *class) String() string {
	return c.name
}

func (c *class) arity() int {
	if initializer, ok := c.methods["init"]; ok {
		return initializer.arity()
	}
	return 0
}

func (c *class) call(interpreter *Interpreter, args []any) (any, error) {
	instance := &instance{c, make(map[string]any)}
	initializer, ok := c.methods["init"]
	if ok {
		_, err := initializer.bind(instance).call(interpreter, args)
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

type instance struct {
	class  *class
	fields map[string]any
}

func (i *instance) get(name ast.Token) (any, error) {
	value, ok := i.fields[name.Lexeme]
	if !ok {
		if method, ok := i.class.methods[name.Lexeme]; ok {
			return method.bind(i), nil
		}
		return nil, &Error{name, ErrUndefinedProperty}
	}
	return value, nil
}

func (i *instance) set(name ast.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *instance) String() string {
	return fmt.Sprintf("%s instance", i.class)
}
