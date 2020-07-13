package parser

import (
	"testing"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/toks"
	"github.com/maleksiuk/golox/tools"
)

func assertAST(t *testing.T, expression expr.Expr, expectedTree string) {
	actualTree := tools.PrintAst(expression)
	if actualTree != expectedTree {
		t.Errorf("Expected AST to be %v but it was %v", expectedTree, actualTree)
	}
}

func TestParseArithmeticAndComparison(t *testing.T) {
	// (123.9 + 92) >= 5 * -9
	tokens := []toks.Token{
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "123.9", Literal: 123.9, Line: 0},
		{TokenType: toks.Plus, Lexeme: "+", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "92", Literal: 92, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.GreaterEqual, Lexeme: ">=", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "5", Literal: 5, Line: 0},
		{TokenType: toks.Star, Lexeme: "*", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "-9", Literal: -9, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	expression := Parse(tokens, &errorreport.ErrorReport{})
	assertAST(t, expression, "(>= (group (+ 123.9 92)) (* 5 -9))")
}

func TestParseUnariesStringsAndBooleans(t *testing.T) {
	// !("str1" == "str2") == false
	tokens := []toks.Token{
		{TokenType: toks.Bang, Lexeme: "!", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"str1\"", Literal: "str1", Line: 0},
		{TokenType: toks.EqualEqual, Lexeme: "==", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"str2\"", Literal: "str2", Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.EqualEqual, Lexeme: "==", Literal: nil, Line: 0},
		{TokenType: toks.False, Lexeme: "false", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	expression := Parse(tokens, &errorreport.ErrorReport{})
	assertAST(t, expression, "(== (! (group (== str1 str2))) false)")
}
