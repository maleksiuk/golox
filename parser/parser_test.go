package parser

import (
	"testing"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/stmt"
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
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	statements := Parse(tokens, &errorreport.ErrorReport{})
	expression := statements[0].(*stmt.Expression).Expression

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
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	statements := Parse(tokens, &errorreport.ErrorReport{})
	expression := statements[0].(*stmt.Expression).Expression
	assertAST(t, expression, "(== (! (group (== str1 str2))) false)")
}

func TestParseVariableDeclarations(t *testing.T) {
	tokens := []toks.Token{
		{TokenType: toks.Var, Lexeme: "var", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "hello", Literal: nil, Line: 0},
		{TokenType: toks.Equal, Lexeme: "=", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "55", Literal: 55, Line: 0},
		{TokenType: toks.Plus, Lexeme: "+", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "33", Literal: 33, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	statements := Parse(tokens, &errorreport.ErrorReport{})
	nameToken := statements[0].(*stmt.Var).Name
	if nameToken.Lexeme != "hello" {
		t.Errorf("Expected name token's lexeme to be 'hello'")
	}

	initializerExpr := statements[0].(*stmt.Var).Initializer
	assertAST(t, initializerExpr, "(+ 55 33)")
}

func TestParseVariableAssignments(t *testing.T) {
	tokens := []toks.Token{
		{TokenType: toks.Identifier, Lexeme: "hello", Literal: nil, Line: 0},
		{TokenType: toks.Equal, Lexeme: "=", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "55", Literal: 55, Line: 0},
		{TokenType: toks.Plus, Lexeme: "+", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "33", Literal: 33, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	statements := Parse(tokens, &errorreport.ErrorReport{})
	expression := statements[0].(*stmt.Expression).Expression

	assertAST(t, expression, "(= hello (+ 55 33))")
}
