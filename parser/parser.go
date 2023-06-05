package parser

import (
	"github.com/lukibw/glox/ast"
)

type Parser struct {
	tokens  []ast.Token
	current int
}

func New(tokens []ast.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) newError(kind ErrorKind) error {
	return &Error{p.peek(), kind}
}

func (p *Parser) peek() ast.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() ast.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == ast.Eof
}

func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(kind ast.TokenKind) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Kind == kind
}

func (p *Parser) match(kinds ...ast.TokenKind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenKind ast.TokenKind, errorKind ErrorKind) (ast.Token, error) {
	if p.check(tokenKind) {
		return p.advance(), nil
	}
	return ast.Token{}, p.newError(errorKind)
}

func (p *Parser) Run() ([]ast.Stmt, []error) {
	statements := make([]ast.Stmt, 0)
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

var newStatementTokenKinds = []ast.TokenKind{
	ast.Class,
	ast.Fun,
	ast.Var,
	ast.For,
	ast.If,
	ast.While,
	ast.Print,
	ast.Return,
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Kind == ast.Semicolon {
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

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(ast.Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(ast.Identifier, ErrMissingVariableName)
	if err != nil {
		return nil, err
	}
	var initializer ast.Expr
	if p.match(ast.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(ast.Semicolon, ErrMissingVarSemicolon)
	if err != nil {
		return nil, err
	}
	return &ast.VarStmt{Name: name, Initializer: initializer}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	switch {
	case p.match(ast.Print):
		return p.printStatement()
	case p.match(ast.LeftBrace):
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.BlockStmt{Statements: statements}, nil
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(ast.Semicolon, ErrMissingValueSemicolon); err != nil {
		return nil, err
	}
	return &ast.PrintStmt{Expression: expr}, nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := make([]ast.Stmt, 0)
	for !p.check(ast.RightBrace) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}
	if _, err := p.consume(ast.RightBrace, ErrMissingRightBrace); err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(ast.Semicolon, ErrMissingExprSemicolon); err != nil {
		return nil, err
	}
	return &ast.ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match(ast.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if v, ok := expr.(*ast.VarExpr); ok {
			return &ast.AssignExpr{Name: v.Name, Value: value}, nil
		}
		return nil, &Error{equals, ErrInvalidAssignTarget}
	}
	return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(ast.BangEqual, ast.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(ast.Greater, ast.GreaterEqual, ast.Less, ast.LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(ast.Minus, ast.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(ast.Slash, ast.Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(ast.Bang, ast.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	switch {
	case p.match(ast.False):
		return &ast.LiteralExpr{Value: false}, nil
	case p.match(ast.True):
		return &ast.LiteralExpr{Value: true}, nil
	case p.match(ast.Nil):
		return &ast.LiteralExpr{Value: nil}, nil
	case p.match(ast.Number, ast.String):
		return &ast.LiteralExpr{Value: p.previous().Literal}, nil
	case p.match(ast.Identifier):
		return &ast.VarExpr{Name: p.previous()}, nil
	case p.match(ast.LeftParen):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err = p.consume(ast.RightParen, ErrMissingRightParen); err != nil {
			return nil, err
		}
		return &ast.GroupingExpr{Expression: expr}, nil
	}
	return nil, p.newError(ErrMissingExpr)
}
