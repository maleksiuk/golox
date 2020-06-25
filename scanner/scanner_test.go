package scanner

import (
	"testing"

	"github.com/maleksiuk/golox/tokens"
)

func TestScanTokens(t *testing.T) {
	tokenSlice := ScanTokens("()")

	if len(tokenSlice) != 3 {
		t.Errorf("len(tokens) should be 3, was %d", len(tokenSlice))
	}

	if tokenSlice[0].TokenType != tokens.LeftParen {
		t.Errorf("First token should be of type %v, was %v", tokens.LeftParen, tokenSlice[0].TokenType)
	}

	if tokenSlice[1].TokenType != tokens.RightParen {
		t.Errorf("Second token should be of type %v, was %v", tokens.RightParen, tokenSlice[1].TokenType)
	}

	if tokenSlice[2].TokenType != tokens.EOF {
		t.Errorf("Last token should be of type %v, was %v", tokens.EOF, tokenSlice[2].TokenType)
	}

}
