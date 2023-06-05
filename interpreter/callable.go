package interpreter

import (
	"fmt"
	"time"

	"github.com/lukibw/glox/ast"
)

type Callable interface {
	Arity() int
	Call(*Interpreter, []any) (any, error)
}

type Function struct {
	Declaration *ast.FunctionStmt
	Closure     *env
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f *Function) Call(interpreter *Interpreter, args []any) (any, error) {
	env := newEnv(f.Closure)
	for i := 0; i < len(f.Declaration.Params); i++ {
		env.define(f.Declaration.Params[i].Lexeme, args[i])
	}
	err := interpreter.executeBlock(f.Declaration.Body, env)
	if rerr, ok := err.(*returnError); ok {
		return rerr.value, nil
	}
	return nil, err
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}

type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(interpreter *Interpreter, args []any) (any, error) {
	return float64(time.Now().Unix()) / 1000.0, nil
}

func (c *Clock) String() string {
	return "<native fn>"
}
