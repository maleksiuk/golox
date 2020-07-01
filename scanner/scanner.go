package scanner

import (
	"strconv"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/srccode"
	"github.com/maleksiuk/golox/toks"
)

var keywords = map[string]toks.TokenType{
	"and":    toks.And,
	"class":  toks.Class,
	"else":   toks.Else,
	"false":  toks.False,
	"for":    toks.For,
	"fun":    toks.Fun,
	"if":     toks.If,
	"nil":    toks.Nil,
	"or":     toks.Or,
	"print":  toks.Print,
	"return": toks.Return,
	"super":  toks.Super,
	"this":   toks.This,
	"true":   toks.True,
	"var":    toks.Var,
	"while":  toks.While,
}

// ScanTokens extracts tokens from a string of Lox code
func ScanTokens(sourceStr string, errorReport *errorreport.ErrorReport) []toks.Token {
	source := srccode.NewSource(sourceStr)

	// our number of tokens will probably be less than half the source length, so we could revise this later
	tokens := make([]toks.Token, 0, source.Len()/2)

	for !source.AtEnd() {
		source.BeginNewLexeme()
		scanToken(&source, &tokens, errorReport)
	}

	addToken(&tokens, toks.EOF, nil)

	return tokens
}

func scanToken(source *srccode.Source, tokens *[]toks.Token, errorReport *errorreport.ErrorReport) {
	r := source.Advance()

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
		if source.Match('=') {
			addToken(tokens, toks.BangEqual, nil)
		} else {
			addToken(tokens, toks.Bang, nil)
		}
	case '=':
		if source.Match('=') {
			addToken(tokens, toks.EqualEqual, nil)
		} else {
			addToken(tokens, toks.Equal, nil)
		}
	case '<':
		if source.Match('=') {
			addToken(tokens, toks.LessEqual, nil)
		} else {
			addToken(tokens, toks.Less, nil)
		}
	case '>':
		if source.Match('=') {
			addToken(tokens, toks.GreaterEqual, nil)
		} else {
			addToken(tokens, toks.Greater, nil)
		}
	case '/':
		if source.Match('/') {
			// ignore commented line
			for source.Peek() != '\n' && !source.AtEnd() {
				source.Advance()
			}
		} else {
			addToken(tokens, toks.Slash, nil)
		}
	case '"':
		handleString(source, tokens, errorReport)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		handleNumber(source, tokens, errorReport)
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		source.IncrementLine()
	default:
		if isAlpha(r) {
			handleIdentifier(source, tokens, errorReport)
		} else {
			errorReport.Report(source.CurrentLine(), "", "Unexpected character.")
		}
	}
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func handleIdentifier(source *srccode.Source, tokens *[]toks.Token, errorReport *errorreport.ErrorReport) {
	for isAlphaNumeric(source.Peek()) {
		source.Advance()
	}

	text := source.Substring(0, 0)
	tokenType, keyExists := keywords[text]
	if keyExists {
		addToken(tokens, tokenType, nil)
	} else {
		addToken(tokens, toks.Identifier, nil)
	}
}

func handleNumber(source *srccode.Source, tokens *[]toks.Token, errorReport *errorreport.ErrorReport) {
	for isDigit(source.Peek()) {
		source.Advance()
	}

	// Look for a fractional part.
	if source.Peek() == '.' && isDigit(source.PeekNext()) {
		// Consume the "."
		source.Advance()

		for isDigit(source.Peek()) {
			source.Advance()
		}
	}

	numStr := source.Substring(0, 0)
	numValue, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		errorReport.Report(source.CurrentLine(), "", "Could not convert number literal to float.")
		return
	}

	addToken(tokens, toks.Number, numValue)
}

func handleString(source *srccode.Source, tokens *[]toks.Token, errorReport *errorreport.ErrorReport) {
	for source.Peek() != '"' && !source.AtEnd() {
		if source.Peek() == '\n' {
			source.IncrementLine()
		}
		source.Advance()
	}

	// Unterminated string.
	if source.AtEnd() {
		errorReport.Report(source.CurrentLine(), "", "Unterminated string.")
		return
	}

	// The closing ".
	source.Advance()

	// Trim the surrounding quotes.
	strValue := source.Substring(1, -1)
	addToken(tokens, toks.String, strValue)
}

func addToken(tokens *[]toks.Token, tokenType toks.TokenType, value interface{}) {
	*tokens = append(*tokens, toks.Token{TokenType: tokenType, Literal: value})
}
