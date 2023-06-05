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

func (p *Parser[T]) consume(tokenKind token.Kind, errorKind ErrorKind) (token.Token, error) {
	if p.check(tokenKind) {
		return p.advance(), nil
	}
	return token.Token{}, p.newError(errorKind)
}

func (p *Parser[T]) Parse() ([]ast.Stmt[T], []error) {
	statements := make([]ast.Stmt[T], 0)
	errors := make([]error, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errors = append(errors, err)
			p.synchronize()
		} else {
			statements = append(statements, stmt)
		}
	}
	return statements, errors
}

func (p *Parser[T]) declaration() (ast.Stmt[T], error) {
	if p.match(token.Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser[T]) varDeclaration() (ast.Stmt[T], error) {
	name, err := p.consume(token.Identifier, ErrMissingVariableName)
	if err != nil {
		return nil, err
	}
	var initializer ast.Expr[T]
	if p.match(token.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.Semicolon, ErrMissingVarSemicolon)
	if err != nil {
		return nil, err
	}
	return &ast.VarStmt[T]{Name: name, Initializer: initializer}, nil
}

func (p *Parser[T]) statement() (ast.Stmt[T], error) {
	switch {
	case p.match(token.Print):
		return p.printStatement()
	case p.match(token.LeftBrace):
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.BlockStmt[T]{Statements: statements}, nil
	default:
		return p.expressionStatement()
	}
}

func (p *Parser[T]) printStatement() (ast.Stmt[T], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.Semicolon, ErrMissingValueSemicolon); err != nil {
		return nil, err
	}
	return &ast.PrintStmt[T]{Expression: expr}, nil
}

func (p *Parser[T]) block() ([]ast.Stmt[T], error) {
	statements := make([]ast.Stmt[T], 0)
	for !p.check(token.RightBrace) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}
	if _, err := p.consume(token.RightBrace, ErrMissingRightBrace); err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser[T]) expressionStatement() (ast.Stmt[T], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.Semicolon, ErrMissingExprSemicolon); err != nil {
		return nil, err
	}
	return &ast.ExpressionStmt[T]{Expression: expr}, nil
}

func (p *Parser[T]) expression() (ast.Expr[T], error) {
	return p.equality()
}

func (p *Parser[T]) assignment() (ast.Expr[T], error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match(token.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if v, ok := expr.(*ast.VariableExpr[T]); ok {
			return &ast.AssignExpr[T]{Name: v.Name, Value: value}, nil
		}
		return nil, &Error{equals, ErrAssignTarget}
	}
	return expr, nil
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
	case p.match(token.Identifier):
		return &ast.VariableExpr[T]{Name: p.previous()}, nil
	case p.match(token.LeftParen):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err = p.consume(token.RightParen, ErrMissingRightParen); err != nil {
			return nil, err
		}
		return &ast.GroupingExpr[T]{Expression: expr}, nil
	}
	return nil, p.newError(ErrMissingExpr)
}

var newStatementTokenKinds = []token.Kind{
	token.Class,
	token.Fun,
	token.Var,
	token.For,
	token.If,
	token.While,
	token.Print,
	token.Return,
}

func (p *Parser[T]) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Kind == token.Semicolon {
			return
		}
		kind := p.peek().Kind
		for _, k := range newStatementTokenKinds {
			if kind == k {
				return
			}
		}
		p.advance()
	}
}
