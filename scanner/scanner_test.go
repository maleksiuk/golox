package scanner

import (
	"testing"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/tokens"
)

func TestScanTokens(t *testing.T) {
	tokenSlice := ScanTokens("()", &errorreport.ErrorReport{})

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

func TestScanTokensWithMultipleCharacters(t *testing.T) {
	bangTokenSlice := ScanTokens("!", &errorreport.ErrorReport{})
	bangEqualTokenSlice := ScanTokens("!=", &errorreport.ErrorReport{})

	if len(bangTokenSlice) != 2 {
		t.Errorf("len(bangTokenSlice) should be 2, was %d", len(bangTokenSlice))
	}

	if len(bangEqualTokenSlice) != 2 {
		t.Errorf("len(bangEqualTokenSlice) should be 2, was %d", len(bangEqualTokenSlice))
	}

	if bangTokenSlice[0].TokenType != tokens.Bang {
		t.Errorf("bangTokenSlice token should be of type %v, was %v", tokens.Bang, bangTokenSlice[0].TokenType)
	}

	if bangEqualTokenSlice[0].TokenType != tokens.BangEqual {
		t.Errorf("bangEqualTokenSlice token should be of type %v, was %v", tokens.BangEqual, bangEqualTokenSlice[0].TokenType)
	}
}
