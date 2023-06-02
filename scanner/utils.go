package scanner

import "github.com/lukibw/glox/token"

var keywords = map[string]token.Kind{
	"and":    token.And,
	"class":  token.Class,
	"else":   token.Else,
	"false":  token.False,
	"for":    token.For,
	"fun":    token.Fun,
	"if":     token.If,
	"nil":    token.Nil,
	"or":     token.Or,
	"print":  token.Print,
	"return": token.Return,
	"super":  token.Super,
	"this":   token.This,
	"true":   token.True,
	"var":    token.Var,
	"while":  token.While,
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isDigit(r) || isAlpha(r)
}
