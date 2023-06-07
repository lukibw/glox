package interpreter

import (
	"fmt"
	"time"

	"github.com/lukibw/glox/ast"
)

type callable interface {
	arity() int
	call(*Interpreter, []any) (any, error)
}

type function struct {
	declaration *ast.FunctionStmt
	closure     *env
}

func (f *function) arity() int {
	return len(f.declaration.Params)
}

func (f *function) call(interpreter *Interpreter, args []any) (any, error) {
	env := newEnv(f.closure)
	for i := 0; i < len(f.declaration.Params); i++ {
		env.define(f.declaration.Params[i].Lexeme, args[i])
	}
	return interpreter.executeBlock(f.declaration.Body, env)
}

func (f *function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

type clock struct{}

func (c *clock) arity() int {
	return 0
}

func (c *clock) call(interpreter *Interpreter, args []any) (any, error) {
	return float64(time.Now().Unix()) / 1000.0, nil
}

func (c *clock) String() string {
	return "<native fn>"
}
