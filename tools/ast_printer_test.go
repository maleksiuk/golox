package tools

import (
	"testing"

	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/toks"
)

func TestPrintAst(t *testing.T) {
	minus := toks.Token{TokenType: toks.Minus, Lexeme: "-"}
	star := toks.Token{TokenType: toks.Star, Lexeme: "*"}
	grouping := expr.Grouping{Expression: &expr.Literal{Value: 45.67}}
	unary := expr.Unary{Operator: minus, Right: &expr.Literal{Value: 123}}
	binary := expr.Binary{Left: &unary, Operator: star, Right: &grouping}

	str := PrintAst(&binary)

	expected := "(* (- 123) (group 45.67))"
	if str != expected {
		t.Errorf("AstPrinter.Print() = %v, want %v", str, expected)
	}
}
