package scanner

import (
	"strconv"

	"github.com/lukibw/glox/ast"
)

type Scanner struct {
	source  []rune
	start   int
	current int
	line    int
	tokens  []ast.Token
}

func New(source string) *Scanner {
	return &Scanner{[]rune(source), 0, 0, 1, make([]ast.Token, 0)}
}

func (s *Scanner) newError(kind ErrorKind) error {
	return &Error{s.line, kind}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) match(r rune) bool {
	if s.isAtEnd() || s.source[s.current] != r {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return -1
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return -1
	}
	return s.source[s.current+1]
}

func (s *Scanner) addTokenWithLiteral(kind ast.TokenKind, literal interface{}) {
	s.tokens = append(
		s.tokens,
		ast.Token{
			Kind:    kind,
			Lexeme:  string(s.source[s.start:s.current]),
			Literal: literal,
			Line:    s.line,
		},
	)
}

func (s *Scanner) addToken(kind ast.TokenKind) {
	s.addTokenWithLiteral(kind, nil)
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return s.newError(ErrUnterminatedString)
	}
	s.advance()
	value := string(s.source[(s.start + 1):(s.current - 1)])
	s.addTokenWithLiteral(ast.String, value)
	return nil
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	value, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		panic("scanner: cannot parse a float")
	}
	s.addTokenWithLiteral(ast.Number, value)
}

var keywords = map[string]ast.TokenKind{
	"and":    ast.And,
	"class":  ast.Class,
	"else":   ast.Else,
	"false":  ast.False,
	"for":    ast.For,
	"fun":    ast.Fun,
	"if":     ast.If,
	"nil":    ast.Nil,
	"or":     ast.Or,
	"print":  ast.Print,
	"return": ast.Return,
	"super":  ast.Super,
	"this":   ast.This,
	"true":   ast.True,
	"var":    ast.Var,
	"while":  ast.While,
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	value := string(s.source[s.start:s.current])
	kind, ok := keywords[value]
	if !ok {
		kind = ast.Identifier
	}
	s.addToken(kind)
}

func (s *Scanner) token() error {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(ast.LeftParen)
	case ')':
		s.addToken(ast.RightParen)
	case '{':
		s.addToken(ast.LeftBrace)
	case '}':
		s.addToken(ast.RightBrace)
	case ',':
		s.addToken(ast.Comma)
	case '.':
		s.addToken(ast.Dot)
	case '-':
		s.addToken(ast.Minus)
	case '+':
		s.addToken(ast.Plus)
	case ';':
		s.addToken(ast.Semicolon)
	case '*':
		s.addToken(ast.Star)
	case '!':
		if s.match('=') {
			s.addToken(ast.BangEqual)
		} else {
			s.addToken(ast.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EqualEqual)
		} else {
			s.addToken(ast.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LessEqual)
		} else {
			s.addToken(ast.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GreaterEqual)
		} else {
			s.addToken(ast.Greater)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(ast.Slash)
		}
	case ' ', '\t', '\r':
		break
	case '\n':
		s.line++
	case '"':
		if err := s.string(); err != nil {
			return err
		}
	default:
		switch {
		case isDigit(r):
			s.number()
		case isAlpha(r):
			s.identifier()
		default:
			return s.newError(ErrUnexpectedCharacter)
		}
	}
	return nil
}

func (s *Scanner) Run() ([]ast.Token, error) {
	var err error
	for !s.isAtEnd() {
		s.start = s.current
		if err = s.token(); err != nil {
			return nil, err
		}
	}
	s.tokens = append(
		s.tokens,
		ast.Token{Kind: ast.Eof, Line: s.line},
	)
	return s.tokens, nil
}
