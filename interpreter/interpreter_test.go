package interpreter

import (
	"testing"
	"time"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/parser"
	"github.com/maleksiuk/golox/scanner"
	"github.com/maleksiuk/golox/stmt"
)

func newMockErrorReport() errorreport.ErrorReport {
	return errorreport.ErrorReport{Printer: errorreport.NewMockPrinter()}
}

func scanAndParse(code string) []stmt.Stmt {
	errorReport := newMockErrorReport()

	tokens := scanner.ScanTokens(code, &errorReport)
	return parser.Parse(tokens, &errorReport)
}

func TestInterpretOrStatement(t *testing.T) {
	code := `
	  var a = 5;
	  var b = 10;
	  var shouldBeTrue = a + b < 3 or a + b > 13;
	  var shouldBeFalse = a + b < 3 or a + b > 18;
	`
	statements := scanAndParse(code)

	errorReport := newMockErrorReport()
	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorReport)
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
	statements := scanAndParse(code)

	errorReport := newMockErrorReport()
	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorReport)
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
	statements := scanAndParse(code)

	errorReport := newMockErrorReport()
	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorReport)
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
	statements := scanAndParse(code)

	errorReport := newMockErrorReport()
	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorReport)
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
	statements := scanAndParse(code)

	errorReport := newMockErrorReport()
	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorReport)
	result := interpreter.GetVariableValue("result").(float64)

	var expected = 4.0
	if result != expected {
		t.Errorf("Expected result to be %v, but it was %v.", expected, result)
	}
}

func TestClockFunction(t *testing.T) {
	code := `
	var result = clock();
  `
	statements := scanAndParse(code)

	secondsSinceEpoch := float64(time.Now().Unix())

	errorReport := newMockErrorReport()
	interpreter := NewInterpreter()
	interpreter.Interpret(statements, &errorReport)
	result := interpreter.GetVariableValue("result").(float64)

	timeDiff := result - secondsSinceEpoch

	if timeDiff > 1 {
		t.Errorf("Expected result to be within one second of %v, but it was not. Result is %v.", secondsSinceEpoch, result)
	}
}
