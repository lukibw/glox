package scan

import (
	"fmt"
	"strconv"
)

type Lexer struct {
	source  []rune
	start   int
	current int
	line    int
	tokens  []Token
}

func NewLexer(source string) *Lexer {
	return &Lexer{[]rune(source), 0, 0, 1, make([]Token, 0)}
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() rune {
	r := l.source[l.current]
	l.current++
	return r
}

func (l *Lexer) match(r rune) bool {
	if l.isAtEnd() || l.source[l.current] != r {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return -1
	}
	return l.source[l.current]
}

func (l *Lexer) peekNext() rune {
	if l.current+1 >= len(l.source) {
		return -1
	}
	return l.source[l.current+1]
}

func (l *Lexer) addTokenWithLiteral(kind TokenKind, literal interface{}) {
	l.tokens = append(l.tokens, Token{kind, string(l.source[l.start:l.current]), literal, l.line})
}

func (l *Lexer) addToken(kind TokenKind) {
	l.addTokenWithLiteral(kind, nil)
}

func (l *Lexer) string() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
	if l.isAtEnd() {
		fmt.Println("Unterminated string")
		return
	}
	l.advance()
	value := string(l.source[(l.start + 1):(l.current - 1)])
	l.addTokenWithLiteral(String, value)
}

func (l *Lexer) number() {
	for isDigit(l.peek()) {
		l.advance()
	}
	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
	}
	value, err := strconv.ParseFloat(string(l.source[l.start:l.current]), 64)
	if err != nil {
		panic("scan: cannot parse a float")
	}
	l.addTokenWithLiteral(Number, value)
}

func (l *Lexer) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}
	value := string(l.source[l.start:l.current])
	kind, ok := keywords[value]
	if !ok {
		kind = Identifier
	}
	l.addToken(kind)
}

func (l *Lexer) scanToken() {
	r := l.advance()
	switch r {
	case '(':
		l.addToken(LeftParen)
	case ')':
		l.addToken(RightParen)
	case '{':
		l.addToken(LeftBrace)
	case '}':
		l.addToken(RightBrace)
	case ',':
		l.addToken(Comma)
	case '.':
		l.addToken(Dot)
	case '-':
		l.addToken(Minus)
	case '+':
		l.addToken(Plus)
	case ';':
		l.addToken(Semicolon)
	case '*':
		l.addToken(Star)
	case '!':
		if l.match('=') {
			l.addToken(BangEqual)
		} else {
			l.addToken(Bang)
		}
	case '=':
		if l.match('=') {
			l.addToken(EqualEqual)
		} else {
			l.addToken(Equal)
		}
	case '<':
		if l.match('=') {
			l.addToken(LessEqual)
		} else {
			l.addToken(Less)
		}
	case '>':
		if l.match('=') {
			l.addToken(GreaterEqual)
		} else {
			l.addToken(Greater)
		}
	case '/':
		if l.match('/') {
			// A comment goes until the end of the line.
			for !l.isAtEnd() && l.peek() != '\n' {
				l.advance()
			}
		} else {
			l.addToken(Slash)
		}
	case ' ', '\t', '\r':
		break
	case '\n':
		l.line++
	case '"':
		l.string()
	default:
		switch {
		case isDigit(r):
			l.number()
		case isAlpha(r):
			l.identifier()
		default:
			fmt.Printf("Unexpected character in line %d\n", l.line)
		}
	}
}

func (l *Lexer) ScanTokens() ([]Token, error) {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	l.tokens = append(l.tokens, Token{Eof, "", nil, l.line})
	return l.tokens, nil
}
