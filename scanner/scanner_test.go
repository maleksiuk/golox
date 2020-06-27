package scanner

import (
	"testing"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/toks"
)

func assertSliceLength(t *testing.T, tokens []toks.Token, expectedLength int) {
	if len(tokens) != expectedLength {
		t.Errorf("Slice length should be %d, was %d", expectedLength, len(tokens))
	}
}

func assertTokenType(t *testing.T, token toks.Token, expectedTokenType toks.TokenType) {
	if token.TokenType != expectedTokenType {
		t.Errorf("Token should be of type %v, was %v", expectedTokenType, token.TokenType)
	}
}

func TestScanTokens(t *testing.T) {
	tokens := ScanTokens("()", &errorreport.ErrorReport{})

	assertSliceLength(t, tokens, 3)

	assertTokenType(t, tokens[0], toks.LeftParen)
	assertTokenType(t, tokens[1], toks.RightParen)
	assertTokenType(t, tokens[2], toks.EOF)
}

func TestScanTokensWithMultipleCharacters(t *testing.T) {
	bangTokens := ScanTokens("!", &errorreport.ErrorReport{})
	bangEqualTokens := ScanTokens("!=", &errorreport.ErrorReport{})

	assertSliceLength(t, bangTokens, 2)
	assertSliceLength(t, bangEqualTokens, 2)

	assertTokenType(t, bangTokens[0], toks.Bang)
	assertTokenType(t, bangEqualTokens[0], toks.BangEqual)
}

func TestScanComments(t *testing.T) {
	commentTokens := ScanTokens("// This should be ignored", &errorreport.ErrorReport{})
	slashTokens := ScanTokens("/*", &errorreport.ErrorReport{})

	assertSliceLength(t, commentTokens, 1)
	assertTokenType(t, commentTokens[0], toks.EOF)

	assertSliceLength(t, slashTokens, 3)
	assertTokenType(t, slashTokens[0], toks.Slash)
}

func TestScanMultipleLines(t *testing.T) {
	tokens := ScanTokens("()\n!=", &errorreport.ErrorReport{})

	assertSliceLength(t, tokens, 4)
	assertTokenType(t, tokens[2], toks.BangEqual)
}

func TestScanStrings(t *testing.T) {
	tokens := ScanTokens("\"hello\nthere man\"", &errorreport.ErrorReport{})
	assertSliceLength(t, tokens, 2)
	assertTokenType(t, tokens[0], toks.String)

	if tokens[0].Literal != "hello\nthere man" {
		t.Errorf("Expected token literal to be %q but it was %q", "hello\nthere man", tokens[0].Literal)
	}
}

// TODO: Test that line count is being incremented, including in the comment and string cases.
// TODO: Test unterminated string error.
