package scanner

import (
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
func ScanTokens(source string) []tokens.Token {
	location := sourceLocation{Line: 1}
	runes := []rune(source)

	// our number of tokens will probably be less than half the source length, so we could revise this later
	tokenSlice := make([]tokens.Token, 0, len(runes)/2)

	for !location.atEnd(runes) {
		location.beginNewLexeme()
		scanToken(&location, runes, &tokenSlice)
	}

	addToken(&tokenSlice, tokens.EOF)

	return tokenSlice
}

func scanToken(location *sourceLocation, runes []rune, tokenSlice *[]tokens.Token) {
	r := runes[location.Current]
	location.Current++

	switch r {
	case '(':
		addToken(tokenSlice, tokens.LeftParen)
	case ')':
		addToken(tokenSlice, tokens.RightParen)
	}
}

func addToken(tokenSlice *[]tokens.Token, tokenType tokens.TokenType) {
	*tokenSlice = append(*tokenSlice, tokens.Token{TokenType: tokenType})
}
