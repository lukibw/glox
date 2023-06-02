package scanner

import (
	"fmt"
	"strconv"

	"github.com/lukibw/glox/token"
)

type Scanner struct {
	source  []rune
	start   int
	current int
	line    int
	tokens  []token.Token
}

func New(source string) *Scanner {
	return &Scanner{[]rune(source), 0, 0, 1, make([]token.Token, 0)}
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

func (s *Scanner) addTokenWithLiteral(kind token.Kind, literal interface{}) {
	s.tokens = append(
		s.tokens,
		token.Token{
			Kind:    kind,
			Lexeme:  string(s.source[s.start:s.current]),
			Literal: literal,
			Line:    s.line,
		},
	)
}

func (s *Scanner) addToken(kind token.Kind) {
	s.addTokenWithLiteral(kind, nil)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		fmt.Println("Unterminated string")
		return
	}
	s.advance()
	value := string(s.source[(s.start + 1):(s.current - 1)])
	s.addTokenWithLiteral(token.String, value)
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
		panic("scan: cannot parse a float")
	}
	s.addTokenWithLiteral(token.Number, value)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	value := string(s.source[s.start:s.current])
	kind, ok := keywords[value]
	if !ok {
		kind = token.Identifier
	}
	s.addToken(kind)
}

func (s *Scanner) scanToken() {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(token.LeftParen)
	case ')':
		s.addToken(token.RightParen)
	case '{':
		s.addToken(token.LeftBrace)
	case '}':
		s.addToken(token.RightBrace)
	case ',':
		s.addToken(token.Comma)
	case '.':
		s.addToken(token.Dot)
	case '-':
		s.addToken(token.Minus)
	case '+':
		s.addToken(token.Plus)
	case ';':
		s.addToken(token.Semicolon)
	case '*':
		s.addToken(token.Star)
	case '!':
		if s.match('=') {
			s.addToken(token.BangEqual)
		} else {
			s.addToken(token.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EqualEqual)
		} else {
			s.addToken(token.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LessEqual)
		} else {
			s.addToken(token.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GreaterEqual)
		} else {
			s.addToken(token.Greater)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(token.Slash)
		}
	case ' ', '\t', '\r':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		switch {
		case isDigit(r):
			s.number()
		case isAlpha(r):
			s.identifier()
		default:
			fmt.Printf("Unexpected character in line %d\n", s.line)
		}
	}
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(
		s.tokens,
		token.Token{
			Kind:    token.Eof,
			Lexeme:  "",
			Literal: nil,
			Line:    s.line,
		},
	)
	return s.tokens, nil
}
