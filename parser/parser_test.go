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

func TestParseLogicalOperators(t *testing.T) {
	// hello == 55 or true and false and something
	tokens := []toks.Token{
		{TokenType: toks.Identifier, Lexeme: "hello", Literal: nil, Line: 0},
		{TokenType: toks.EqualEqual, Lexeme: "==", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "55", Literal: 55, Line: 0},
		{TokenType: toks.Or, Lexeme: "or", Literal: nil, Line: 0},
		{TokenType: toks.True, Lexeme: "true", Literal: 33, Line: 0},
		{TokenType: toks.And, Lexeme: "and", Literal: nil, Line: 0},
		{TokenType: toks.False, Lexeme: "false", Literal: 33, Line: 0},
		{TokenType: toks.And, Lexeme: "and", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "something", Literal: nil, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	errorReport := errorreport.ErrorReport{}
	statements := Parse(tokens, &errorReport)
	expression := statements[0].(*stmt.Expression).Expression

	assertAST(t, expression, "(or (== hello 55) (and (and true false) something))")
}

func TestParseWhileStatements(t *testing.T) {
	// while (something == 3) {
	//   print "hi";
	// }
	tokens := []toks.Token{
		{TokenType: toks.While, Lexeme: "while", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "something", Literal: nil, Line: 0},
		{TokenType: toks.EqualEqual, Lexeme: "==", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "3", Literal: 3, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.LeftBrace, Lexeme: "{", Literal: nil, Line: 0},
		{TokenType: toks.Print, Lexeme: "print", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"hi\"", Literal: "hi", Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.RightBrace, Lexeme: "}", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	errorReport := errorreport.ErrorReport{}
	statements := Parse(tokens, &errorReport)
	condition := statements[0].(*stmt.While).Condition
	body := statements[0].(*stmt.While).Body
	blockStatements := body.(*stmt.Block).Statements
	printExpression := blockStatements[0].(*stmt.Print).Expression

	assertAST(t, condition, "(== something 3)")
	assertAST(t, printExpression, "hi")
}

func TestParseForStatements(t *testing.T) {
	// for (var x = 1; x < 3; x = x + 1) {
	//   print "hi";
	// }
	tokens := []toks.Token{
		{TokenType: toks.For, Lexeme: "for", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},

		{TokenType: toks.Var, Lexeme: "var", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "x", Literal: nil, Line: 0},
		{TokenType: toks.Equal, Lexeme: "=", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "1", Literal: 1, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},

		{TokenType: toks.Identifier, Lexeme: "x", Literal: nil, Line: 0},
		{TokenType: toks.Less, Lexeme: "<", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "3", Literal: 3, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},

		{TokenType: toks.Identifier, Lexeme: "x", Literal: nil, Line: 0},
		{TokenType: toks.Equal, Lexeme: "=", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "x", Literal: nil, Line: 0},
		{TokenType: toks.Plus, Lexeme: "+", Literal: nil, Line: 0},
		{TokenType: toks.Number, Lexeme: "1", Literal: 1, Line: 0},

		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.LeftBrace, Lexeme: "{", Literal: nil, Line: 0},
		{TokenType: toks.Print, Lexeme: "print", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"hi\"", Literal: "hi", Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.RightBrace, Lexeme: "}", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	errorReport := errorreport.ErrorReport{}
	statements := Parse(tokens, &errorReport)
	outerBlockStatements := statements[0].(*stmt.Block).Statements
	initializerExpression := outerBlockStatements[0].(*stmt.Var).Initializer

	whileStatement := outerBlockStatements[1].(*stmt.While)
	condition := whileStatement.Condition
	body := whileStatement.Body
	bodyStatements := body.(*stmt.Block).Statements
	userSpecified := bodyStatements[0].(*stmt.Block)
	printExpression := userSpecified.Statements[0].(*stmt.Print).Expression
	incrementExpression := bodyStatements[1].(*stmt.Expression).Expression

	assertAST(t, condition, "(< x 3)")
	assertAST(t, initializerExpression, "1")
	assertAST(t, printExpression, "hi")
	assertAST(t, incrementExpression, "(= x (+ x 1))")
}

func TestParseEmptyForStatements(t *testing.T) {
	// for (;;) {
	//   print "hi";
	// }
	tokens := []toks.Token{
		{TokenType: toks.For, Lexeme: "for", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.LeftBrace, Lexeme: "{", Literal: nil, Line: 0},
		{TokenType: toks.Print, Lexeme: "print", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"hi\"", Literal: "hi", Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.RightBrace, Lexeme: "}", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	errorReport := errorreport.ErrorReport{}
	statements := Parse(tokens, &errorReport)
	outerBlockStatements := statements[0].(*stmt.Block).Statements
	whileStatement := outerBlockStatements[0].(*stmt.While)

	condition := whileStatement.Condition.(*expr.Literal)
	body := whileStatement.Body
	bodyStatements := body.(*stmt.Block).Statements
	userSpecified := bodyStatements[0].(*stmt.Block)
	printExpression := userSpecified.Statements[0].(*stmt.Print).Expression
	if len(bodyStatements) != 1 {
		t.Error("There should only be one body statement because there is no increment expression.")
	}

	if condition.Value != true {
		t.Error("Condition should be a True literal")
	}
	assertAST(t, printExpression, "hi")
}

func TestParseIfStatements(t *testing.T) {
	// if (a != b) {
	//   print "hi";
	// }
	tokens := []toks.Token{
		{TokenType: toks.If, Lexeme: "if", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "a", Literal: nil, Line: 0},
		{TokenType: toks.BangEqual, Lexeme: "!=", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "b", Literal: nil, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.LeftBrace, Lexeme: "{", Literal: nil, Line: 0},
		{TokenType: toks.Print, Lexeme: "print", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"hi\"", Literal: "hi", Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.RightBrace, Lexeme: "}", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	errorReport := errorreport.ErrorReport{}
	statements := Parse(tokens, &errorReport)
	conditional := statements[0].(*stmt.Conditional)
	condition := conditional.Condition.(*expr.Binary)
	thenBlock := conditional.ThenStatement.(*stmt.Block)

	assertAST(t, condition, "(!= a b)")

	printExpression := thenBlock.Statements[0].(*stmt.Print).Expression
	assertAST(t, printExpression, "hi")
}

func TestParseFunctionCalls(t *testing.T) {
	// somefunction()(otherfunction(x + y, z));
	tokens := []toks.Token{
		{TokenType: toks.Identifier, Lexeme: "somefunction", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "otherfunction", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "x", Literal: nil, Line: 0},
		{TokenType: toks.Plus, Lexeme: "+", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "y", Literal: nil, Line: 0},
		{TokenType: toks.Comma, Lexeme: ",", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "z", Literal: nil, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	errorReport := errorreport.ErrorReport{}
	statements := Parse(tokens, &errorReport)
	expression := statements[0].(*stmt.Expression).Expression

	assertAST(t, expression, "(call (call somefunction) (call otherfunction (+ x y),z))")
}

func TestParseFunctionDeclaration(t *testing.T) {
	tokens := []toks.Token{
		{TokenType: toks.Fun, Lexeme: "fun", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "do_something", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.LeftBrace, Lexeme: "{", Literal: nil, Line: 0},
		{TokenType: toks.Print, Lexeme: "print", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"hi\"", Literal: "hi", Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.RightBrace, Lexeme: "}", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	statements := Parse(tokens, &errorreport.ErrorReport{})
	functionStatement := statements[0].(*stmt.Function)
	nameToken := functionStatement.Name
	if nameToken.Lexeme != "do_something" {
		t.Errorf("Expected name token's lexeme to be 'do_something'")
	}

	if len(functionStatement.Params) != 0 {
		t.Errorf("Expected there to be no function parameters")
	}

	bodyStatements := functionStatement.Body
	printExpression := bodyStatements[0].(*stmt.Print).Expression
	assertAST(t, printExpression, "hi")
}

func TestParseFunctionDeclarationParameters(t *testing.T) {
	tokens := []toks.Token{
		{TokenType: toks.Fun, Lexeme: "fun", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "do_something", Literal: nil, Line: 0},
		{TokenType: toks.LeftParen, Lexeme: "(", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "cool_cool_water", Literal: nil, Line: 0},
		{TokenType: toks.Comma, Lexeme: ",", Literal: nil, Line: 0},
		{TokenType: toks.Identifier, Lexeme: "by_marty_robbins", Literal: nil, Line: 0},
		{TokenType: toks.RightParen, Lexeme: ")", Literal: nil, Line: 0},
		{TokenType: toks.LeftBrace, Lexeme: "{", Literal: nil, Line: 0},
		{TokenType: toks.Print, Lexeme: "print", Literal: nil, Line: 0},
		{TokenType: toks.String, Lexeme: "\"hi\"", Literal: "hi", Line: 0},
		{TokenType: toks.Semicolon, Lexeme: ";", Literal: nil, Line: 0},
		{TokenType: toks.RightBrace, Lexeme: "}", Literal: nil, Line: 0},
		{TokenType: toks.EOF, Lexeme: "", Literal: nil, Line: 0},
	}

	statements := Parse(tokens, &errorreport.ErrorReport{})
	functionStatement := statements[0].(*stmt.Function)
	parameters := functionStatement.Params
	if len(parameters) != 2 {
		t.Errorf("Expected there to be 2 function parameters")
	}

	if parameters[0].Lexeme != "cool_cool_water" || parameters[1].Lexeme != "by_marty_robbins" {
		t.Errorf("Expected parameters to be cool_cool_water and by_marty_robbins")
	}
}
