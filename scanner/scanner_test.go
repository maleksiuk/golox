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

func assertTokenLiteral(t *testing.T, token toks.Token, expectedLiteral interface{}) {
	if token.Literal != expectedLiteral {
		t.Errorf("Expected token literal to be %q but it was %q", expectedLiteral, token.Literal)
	}
}

func assertTokenLexeme(t *testing.T, token toks.Token, expectedLexeme string) {
	if token.Lexeme != expectedLexeme {
		t.Errorf("Expected token lexeme to be %v but it was %v", expectedLexeme, token.Lexeme)
	}
}

func assertTokenLine(t *testing.T, token toks.Token, expectedLine int) {
	if token.Line != expectedLine {
		t.Errorf("Expected token %s to have line %d but had line %d", token.Lexeme, expectedLine, token.Line)
	}
}

func TestScanTokens(t *testing.T) {
	tokens := ScanTokens("()", &errorreport.ErrorReport{})

	assertSliceLength(t, tokens, 3)

	assertTokenType(t, tokens[0], toks.LeftParen)
	assertTokenLexeme(t, tokens[0], "(")
	assertTokenType(t, tokens[1], toks.RightParen)
	assertTokenLexeme(t, tokens[1], ")")
	assertTokenType(t, tokens[2], toks.EOF)
}

func TestScanTokensWithMultipleCharacters(t *testing.T) {
	bangTokens := ScanTokens("!", &errorreport.ErrorReport{})
	bangEqualTokens := ScanTokens("!=", &errorreport.ErrorReport{})

	assertSliceLength(t, bangTokens, 2)
	assertSliceLength(t, bangEqualTokens, 2)

	assertTokenType(t, bangTokens[0], toks.Bang)
	assertTokenType(t, bangEqualTokens[0], toks.BangEqual)

	assertTokenLexeme(t, bangTokens[0], "!")
	assertTokenLexeme(t, bangEqualTokens[0], "!=")
}

func TestScanComments(t *testing.T) {
	commentTokens := ScanTokens("// This should be ignored", &errorreport.ErrorReport{})
	slashTokens := ScanTokens("/*", &errorreport.ErrorReport{})

	assertSliceLength(t, commentTokens, 1)
	assertTokenType(t, commentTokens[0], toks.EOF)

	assertSliceLength(t, slashTokens, 3)
	assertTokenType(t, slashTokens[0], toks.Slash)
}

func TestLineIncrementingForComments(t *testing.T) {
	tokens := ScanTokens("// This should be ignored\n2 + 3", &errorreport.ErrorReport{})
	assertSliceLength(t, tokens, 4)
	assertTokenLexeme(t, tokens[0], "2")
	assertTokenLine(t, tokens[0], 2)
}

func TestScanMultipleLines(t *testing.T) {
	tokens := ScanTokens("()\n!=", &errorreport.ErrorReport{})

	assertSliceLength(t, tokens, 4)
	assertTokenType(t, tokens[2], toks.BangEqual)
	assertTokenLine(t, tokens[0], 1)
	assertTokenLine(t, tokens[1], 1)
	assertTokenLine(t, tokens[2], 2)
}

func TestScanStrings(t *testing.T) {
	tokens := ScanTokens("\"hello\nthere man\" 2", &errorreport.ErrorReport{})
	assertSliceLength(t, tokens, 3)
	assertTokenType(t, tokens[0], toks.String)
	assertTokenLiteral(t, tokens[0], "hello\nthere man")
	assertTokenLexeme(t, tokens[0], "\"hello\nthere man\"")

	// When there's a newline in a string we increment the line number, so the string's token's line
	// number is 2 instead of 1. That seems odd to me.
	assertTokenLine(t, tokens[0], 2)
	assertTokenLine(t, tokens[1], 2)
}

func TestScanNumbers(t *testing.T) {
	tokens := ScanTokens("123 456.78", &errorreport.ErrorReport{})
	assertSliceLength(t, tokens, 3)
	assertTokenType(t, tokens[0], toks.Number)
	assertTokenType(t, tokens[1], toks.Number)
	assertTokenLiteral(t, tokens[0], 123.0)
	assertTokenLiteral(t, tokens[1], 456.78)
}

func TestScanIdentifiers(t *testing.T) {
	tokens := ScanTokens("orchid", &errorreport.ErrorReport{})
	assertSliceLength(t, tokens, 2)
	assertTokenType(t, tokens[0], toks.Identifier)
}

func TestScanKeywords(t *testing.T) {
	tokens := ScanTokens("or and", &errorreport.ErrorReport{})
	assertSliceLength(t, tokens, 3)
	assertTokenType(t, tokens[0], toks.Or)
	assertTokenType(t, tokens[1], toks.And)
}

func TestUnterminatedStringError(t *testing.T) {
	errorReport := errorreport.ErrorReport{}
	ScanTokens("\"hey man", &errorReport)
	if !errorReport.HadError {
		t.Error("Expected error report to think it had an error.")
	}
	if errorReport.HadRuntimeError {
		t.Error("Expected error report to not say it had a runtime error.")
	}
}
