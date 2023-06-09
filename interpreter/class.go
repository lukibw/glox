package interpreter

import (
	"fmt"

	"github.com/lukibw/glox/ast"
)

type class struct {
	name       string
	methods    map[string]*function
	superclass *class
}

func (c *class) String() string {
	return c.name
}

func (c *class) findMethod(name string) *function {
	if method, ok := c.methods[name]; ok {
		return method
	}
	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}
	return nil
}

func (c *class) arity() int {
	if initializer := c.findMethod("init"); initializer != nil {
		return initializer.arity()
	}
	return 0
}

func (c *class) call(interpreter *Interpreter, args []any) (any, error) {
	instance := &instance{c, make(map[string]any)}
	initializer := c.findMethod("init")
	if initializer != nil {
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
		if method := i.class.findMethod(name.Lexeme); method != nil {
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
