package parser

import (
	"github.com/lukibw/glox/expr"
	"github.com/lukibw/glox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) newError(kind ErrorKind) error {
	return &Error{p.peek(), kind}
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == token.Eof
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(kind token.Kind) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Kind == kind
}

func (p *Parser) match(kinds ...token.Kind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) Parse() (expr.Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (expr.Expr, error) {
	exp, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		exp = &expr.Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

func (p *Parser) comparison() (expr.Expr, error) {
	exp, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		exp = &expr.Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

func (p *Parser) term() (expr.Expr, error) {
	exp, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		exp = &expr.Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

func (p *Parser) factor() (expr.Expr, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		exp = &expr.Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

func (p *Parser) unary() (expr.Expr, error) {
	if p.match(token.Bang, token.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &expr.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (expr.Expr, error) {
	switch {
	case p.match(token.False):
		return &expr.Literal{Value: false}, nil
	case p.match(token.True):
		return &expr.Literal{Value: true}, nil
	case p.match(token.Nil):
		return &expr.Literal{Value: nil}, nil
	case p.match(token.Number, token.String):
		return &expr.Literal{Value: p.previous().Literal}, nil
	case p.match(token.LeftParen):
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.check(token.RightParen) {
			p.advance()
		} else {
			return nil, p.newError(ErrMissingRightParen)
		}
		return &expr.Grouping{Expression: exp}, nil
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

// func (p *Parser) synchronize() {
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
