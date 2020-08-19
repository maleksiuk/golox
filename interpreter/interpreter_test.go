package interpreter

import (
	"testing"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/parser"
	"github.com/maleksiuk/golox/scanner"
)

func TestInterpretOrStatement(t *testing.T) {
	code := `
	  var a = 5;
	  var b = 10;
	  var shouldBeTrue = a + b < 3 or a + b > 13;
	  var shouldBeFalse = a + b < 3 or a + b > 18;
	`
	tokens := scanner.ScanTokens(code, &errorreport.ErrorReport{})
	statements := parser.Parse(tokens, &errorreport.ErrorReport{})

	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorreport.ErrorReport{})
	shouldBeTrue := interpreter.GetVariableValue("shouldBeTrue").(bool)
	shouldBeFalse := interpreter.GetVariableValue("shouldBeFalse").(bool)

	if !shouldBeTrue {
		t.Error("Expected shouldBeTrue to be true, but it was not.")
	}

	if shouldBeFalse {
		t.Error("Expected shouldBeFalse to be false, but it was not.")
	}
}

func TestInterpretAndStatement(t *testing.T) {
	code := `
	  var a = 1;
	  var b = 2;
	  var shouldBeTrue = a + b >= 0 and b > -10;
	  var shouldBeFalse = a + b >= 3 and b > 100;
	`
	tokens := scanner.ScanTokens(code, &errorreport.ErrorReport{})
	statements := parser.Parse(tokens, &errorreport.ErrorReport{})

	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorreport.ErrorReport{})
	shouldBeTrue := interpreter.GetVariableValue("shouldBeTrue").(bool)
	shouldBeFalse := interpreter.GetVariableValue("shouldBeFalse").(bool)

	if !shouldBeTrue {
		t.Error("Expected shouldBeTrue to be true, but it was not.")
	}

	if shouldBeFalse {
		t.Error("Expected shouldBeFalse to be false, but it was not.")
	}
}

func TestInterpretArithmetic(t *testing.T) {
	code := `
	  var result = 1 + 12.6 / 3 * 8;
	`
	tokens := scanner.ScanTokens(code, &errorreport.ErrorReport{})
	statements := parser.Parse(tokens, &errorreport.ErrorReport{})

	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorreport.ErrorReport{})
	result := interpreter.GetVariableValue("result").(float64)

	var expected = 34.6
	if result != expected {
		t.Errorf("Expected result to be %v, but it was %v.", expected, result)
	}
}

func TestInterpretWhile(t *testing.T) {
	code := `
	  var result = 0;
	  while (result < 5) {
		  result = result + 1;
	  }
	`
	tokens := scanner.ScanTokens(code, &errorreport.ErrorReport{})
	statements := parser.Parse(tokens, &errorreport.ErrorReport{})

	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorreport.ErrorReport{})
	result := interpreter.GetVariableValue("result").(float64)

	var expected = 5.0
	if result != expected {
		t.Errorf("Expected result to be %v, but it was %v.", expected, result)
	}
}

func TestInterpretFor(t *testing.T) {
	code := `
	  var result = 0;
	  for (var i = 0; i < 5; i = i + 1) {
		  result = i;
	  }
	`
	tokens := scanner.ScanTokens(code, &errorreport.ErrorReport{})
	statements := parser.Parse(tokens, &errorreport.ErrorReport{})

	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorreport.ErrorReport{})
	result := interpreter.GetVariableValue("result").(float64)

	var expected = 4.0
	if result != expected {
		t.Errorf("Expected result to be %v, but it was %v.", expected, result)
	}
}
