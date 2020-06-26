package scanner

import (
	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/tokens"
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
func ScanTokens(source string, errorReport *errorreport.ErrorReport) []tokens.Token {
	location := sourceLocation{Line: 1}
	runes := []rune(source)

	// our number of tokens will probably be less than half the source length, so we could revise this later
	tokenSlice := make([]tokens.Token, 0, len(runes)/2)

	for !location.atEnd(runes) {
		location.beginNewLexeme()
		scanToken(&location, runes, &tokenSlice, errorReport)
	}

	addToken(&tokenSlice, tokens.EOF)

	return tokenSlice
}

func scanToken(location *sourceLocation, runes []rune, tokenSlice *[]tokens.Token, errorReport *errorreport.ErrorReport) {
	r := runes[location.Current]
	location.Current++

	switch r {
	case '(':
		addToken(tokenSlice, tokens.LeftParen)
	case ')':
		addToken(tokenSlice, tokens.RightParen)
	case '{':
		addToken(tokenSlice, tokens.LeftBrace)
	case '}':
		addToken(tokenSlice, tokens.RightBrace)
	case ',':
		addToken(tokenSlice, tokens.Comma)
	case '.':
		addToken(tokenSlice, tokens.Dot)
	case '-':
		addToken(tokenSlice, tokens.Minus)
	case '+':
		addToken(tokenSlice, tokens.Plus)
	case ';':
		addToken(tokenSlice, tokens.Semicolon)
	case '*':
		addToken(tokenSlice, tokens.Star)
	default:
		errorReport.Report(location.Line, "", "Unexpected character.")
	}
}

func addToken(tokenSlice *[]tokens.Token, tokenType tokens.TokenType) {
	*tokenSlice = append(*tokenSlice, tokens.Token{TokenType: tokenType})
}
