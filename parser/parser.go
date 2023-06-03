package parser

import (
	"github.com/lukibw/glox/ast"
	"github.com/lukibw/glox/token"
)

type Parser[T any] struct {
	tokens  []token.Token
	current int
}

func New[T any](tokens []token.Token) *Parser[T] {
	return &Parser[T]{tokens, 0}
}

func (p *Parser[T]) newError(kind ErrorKind) error {
	return &Error{p.peek(), kind}
}

func (p *Parser[T]) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser[T]) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser[T]) isAtEnd() bool {
	return p.peek().Kind == token.Eof
}

func (p *Parser[T]) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser[T]) check(kind token.Kind) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Kind == kind
}

func (p *Parser[T]) match(kinds ...token.Kind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser[T]) Parse() (ast.Expr[T], error) {
	return p.expression()
}

func (p *Parser[T]) expression() (ast.Expr[T], error) {
	return p.equality()
}

func (p *Parser[T]) equality() (ast.Expr[T], error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser[T]) comparison() (ast.Expr[T], error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser[T]) term() (ast.Expr[T], error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser[T]) factor() (ast.Expr[T], error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser[T]) unary() (ast.Expr[T], error) {
	if p.match(token.Bang, token.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr[T]{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.primary()
}

func (p *Parser[T]) primary() (ast.Expr[T], error) {
	switch {
	case p.match(token.False):
		return &ast.LiteralExpr[T]{Value: false}, nil
	case p.match(token.True):
		return &ast.LiteralExpr[T]{Value: true}, nil
	case p.match(token.Nil):
		return &ast.LiteralExpr[T]{Value: nil}, nil
	case p.match(token.Number, token.String):
		return &ast.LiteralExpr[T]{Value: p.previous().Literal}, nil
	case p.match(token.LeftParen):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.check(token.RightParen) {
			p.advance()
		} else {
			return nil, p.newError(ErrMissingRightParen)
		}
		return &ast.GroupingExpr[T]{Expression: expr}, nil
	}
	return nil, p.newError(ErrMissingExpr)
}

// var newStatementTokenKinds = []token.Kind{
// 	token.Class,
// 	token.Fun,
// 	token.Var,
// 	token.For,
// 	token.If,
// 	token.While,
// 	token.Print,
// 	token.Return,
// }

// func (p *Parser[T]) synchronize() {
// 	p.advance()
// 	for !p.isAtEnd() {
// 		if p.previous().Kind == token.Semicolon {
// 			return
// 		}
// 		kind := p.peek().Kind
// 		for _, k := range newStatementTokenKinds {
// 			if kind == k {
// 				return
// 			}
// 		}
// 		p.advance()
// 	}
// }
