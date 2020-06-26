package scanner

import (
	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/toks"
)

type sourceLocation struct {
	Start   int
	Current int
	Line    int
}

func (location *sourceLocation) atEnd(runes []rune) bool {
	return location.Current >= len(runes)
}

func (location *sourceLocation) beginNewLexeme() {
	location.Start = location.Current
}

// ScanTokens extracts tokens from a string of Lox code
func ScanTokens(source string, errorReport *errorreport.ErrorReport) []toks.Token {
	location := sourceLocation{Line: 1}
	runes := []rune(source)

	// our number of tokens will probably be less than half the source length, so we could revise this later
	tokens := make([]toks.Token, 0, len(runes)/2)

	for !location.atEnd(runes) {
		location.beginNewLexeme()
		scanToken(&location, runes, &tokens, errorReport)
	}

	addToken(&tokens, toks.EOF)

	return tokens
}

func scanToken(location *sourceLocation, runes []rune, tokens *[]toks.Token, errorReport *errorreport.ErrorReport) {
	r := runes[location.Current]
	location.Current++

	switch r {
	case '(':
		addToken(tokens, toks.LeftParen)
	case ')':
		addToken(tokens, toks.RightParen)
	case '{':
		addToken(tokens, toks.LeftBrace)
	case '}':
		addToken(tokens, toks.RightBrace)
	case ',':
		addToken(tokens, toks.Comma)
	case '.':
		addToken(tokens, toks.Dot)
	case '-':
		addToken(tokens, toks.Minus)
	case '+':
		addToken(tokens, toks.Plus)
	case ';':
		addToken(tokens, toks.Semicolon)
	case '*':
		addToken(tokens, toks.Star)
	case '!':
		if match('=', location, runes) {
			addToken(tokens, toks.BangEqual)
		} else {
			addToken(tokens, toks.Bang)
		}
	case '=':
		if match('=', location, runes) {
			addToken(tokens, toks.EqualEqual)
		} else {
			addToken(tokens, toks.Equal)
		}
	case '<':
		if match('=', location, runes) {
			addToken(tokens, toks.LessEqual)
		} else {
			addToken(tokens, toks.Less)
		}
	case '>':
		if match('=', location, runes) {
			addToken(tokens, toks.GreaterEqual)
		} else {
			addToken(tokens, toks.Greater)
		}
	case '/':
		if match('/', location, runes) {
			// ignore commented line
			for peek(location, runes) != '\n' && !location.atEnd(runes) {
				location.Current++
			}
		} else {
			addToken(tokens, toks.Slash)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		location.Line++
	default:
		errorReport.Report(location.Line, "", "Unexpected character.")
	}
}

func addToken(tokens *[]toks.Token, tokenType toks.TokenType) {
	*tokens = append(*tokens, toks.Token{TokenType: tokenType})
}

func match(expected rune, location *sourceLocation, runes []rune) bool {
	if location.atEnd(runes) {
		return false
	}

	if runes[location.Current] != expected {
		return false
	}

	location.Current++
	return true
}

func peek(location *sourceLocation, runes []rune) rune {
	if location.atEnd(runes) {
		return 0
	}

	return runes[location.Current]
}
