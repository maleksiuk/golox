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

func TestParse(t *testing.T) {

	tokens := []toks.Token{
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "123.9", Literal: 123.9, Line: 0},
		{TokenType: toks.Plus, Lexeme: "+", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "92", Literal: 92, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	expression := Parse(tokens, &errorreport.ErrorReport{})
	assertAST(t, expression, "(group (+ 123.9 92))")
}
