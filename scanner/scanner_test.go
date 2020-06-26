package scanner

import (
	"testing"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/tokens"
)

func assertSliceLength(t *testing.T, tokenSlice []tokens.Token, expectedLength int) {
	if len(tokenSlice) != expectedLength {
		t.Errorf("Slice length should be %d, was %d", expectedLength, len(tokenSlice))
	}
}

func assertTokenType(t *testing.T, token tokens.Token, expectedTokenType tokens.TokenType) {
	if token.TokenType != expectedTokenType {
		t.Errorf("Token should be of type %v, was %v", expectedTokenType, token.TokenType)
	}
}

func TestScanTokens(t *testing.T) {
	tokenSlice := ScanTokens("()", &errorreport.ErrorReport{})

	assertSliceLength(t, tokenSlice, 3)

	assertTokenType(t, tokenSlice[0], tokens.LeftParen)
	assertTokenType(t, tokenSlice[1], tokens.RightParen)
	assertTokenType(t, tokenSlice[2], tokens.EOF)
}

func TestScanTokensWithMultipleCharacters(t *testing.T) {
	bangTokenSlice := ScanTokens("!", &errorreport.ErrorReport{})
	bangEqualTokenSlice := ScanTokens("!=", &errorreport.ErrorReport{})

	assertSliceLength(t, bangTokenSlice, 2)
	assertSliceLength(t, bangEqualTokenSlice, 2)

	assertTokenType(t, bangTokenSlice[0], tokens.Bang)
	assertTokenType(t, bangEqualTokenSlice[0], tokens.BangEqual)
}

func TestScanComments(t *testing.T) {
	commentTokenSlice := ScanTokens("// This should be ignored", &errorreport.ErrorReport{})
	slashTokenSlice := ScanTokens("/*", &errorreport.ErrorReport{})

	assertSliceLength(t, commentTokenSlice, 1)
	assertTokenType(t, commentTokenSlice[0], tokens.EOF)

	assertSliceLength(t, slashTokenSlice, 3)
	assertTokenType(t, slashTokenSlice[0], tokens.Slash)
}

func TestScanMultipleLines(t *testing.T) {
	tokenSlice := ScanTokens("()\n!=", &errorreport.ErrorReport{})

	assertSliceLength(t, tokenSlice, 4)
	assertTokenType(t, tokenSlice[2], tokens.BangEqual)
}

// TODO: Test that line count is being incremented, including in the comment case.
