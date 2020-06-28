package scanner

import (
	"strconv"

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

	addToken(&tokens, toks.EOF, nil)

	return tokens
}

func scanToken(location *sourceLocation, runes []rune, tokens *[]toks.Token, errorReport *errorreport.ErrorReport) {
	r := runes[location.Current]
	location.Current++

	switch r {
	case '(':
		addToken(tokens, toks.LeftParen, nil)
	case ')':
		addToken(tokens, toks.RightParen, nil)
	case '{':
		addToken(tokens, toks.LeftBrace, nil)
	case '}':
		addToken(tokens, toks.RightBrace, nil)
	case ',':
		addToken(tokens, toks.Comma, nil)
	case '.':
		addToken(tokens, toks.Dot, nil)
	case '-':
		addToken(tokens, toks.Minus, nil)
	case '+':
		addToken(tokens, toks.Plus, nil)
	case ';':
		addToken(tokens, toks.Semicolon, nil)
	case '*':
		addToken(tokens, toks.Star, nil)
	case '!':
		if match('=', location, runes) {
			addToken(tokens, toks.BangEqual, nil)
		} else {
			addToken(tokens, toks.Bang, nil)
		}
	case '=':
		if match('=', location, runes) {
			addToken(tokens, toks.EqualEqual, nil)
		} else {
			addToken(tokens, toks.Equal, nil)
		}
	case '<':
		if match('=', location, runes) {
			addToken(tokens, toks.LessEqual, nil)
		} else {
			addToken(tokens, toks.Less, nil)
		}
	case '>':
		if match('=', location, runes) {
			addToken(tokens, toks.GreaterEqual, nil)
		} else {
			addToken(tokens, toks.Greater, nil)
		}
	case '/':
		if match('/', location, runes) {
			// ignore commented line
			for peek(location, runes) != '\n' && !location.atEnd(runes) {
				location.Current++
			}
		} else {
			addToken(tokens, toks.Slash, nil)
		}
	case '"':
		handleString(location, tokens, runes, errorReport)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		handleNumber(location, tokens, runes, errorReport)
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		location.Line++
	default:
		errorReport.Report(location.Line, "", "Unexpected character.")
	}
}

func isDigit(r rune) bool {
	return r >= 0x30 && r <= 0x39
}

func handleNumber(location *sourceLocation, tokens *[]toks.Token, runes []rune, errorReport *errorreport.ErrorReport) {
	for isDigit(peek(location, runes)) {
		location.Current++
	}

	// Look for a fractional part.
	if peek(location, runes) == '.' && isDigit(peekNext(location, runes)) {
		// Consume the "."
		location.Current++

		for isDigit(peek(location, runes)) {
			location.Current++
		}
	}

	numStr := string(runes[location.Start:location.Current])
	numValue, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		errorReport.Report(location.Line, "", "Could not convert number literal to float.")
		return
	}

	addToken(tokens, toks.Number, numValue)
}

func handleString(location *sourceLocation, tokens *[]toks.Token, runes []rune, errorReport *errorreport.ErrorReport) {
	for peek(location, runes) != '"' && !location.atEnd(runes) {
		if peek(location, runes) == '\n' {
			location.Line++
		}
		location.Current++
	}

	// Unterminated string.
	if location.atEnd(runes) {
		errorReport.Report(location.Line, "", "Unterminated string.")
		return
	}

	// The closing ".
	location.Current++

	// Trim the surrounding quotes.
	strValue := string(runes[location.Start+1 : location.Current-1])
	addToken(tokens, toks.String, strValue)
}

func addToken(tokens *[]toks.Token, tokenType toks.TokenType, value interface{}) {
	*tokens = append(*tokens, toks.Token{TokenType: tokenType, Literal: value})
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

func peekNext(location *sourceLocation, runes []rune) rune {
	if location.Current+1 >= len(runes) {
		return 0
	}

	return runes[location.Current+1]
}
