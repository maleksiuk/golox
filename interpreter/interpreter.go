package interpreter

import (
	"fmt"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/stmt"
	"github.com/maleksiuk/golox/toks"
)

type interpreter struct {
}

type runtimeError struct {
	token   toks.Token
	message string
}

// Interpret evaluates a program (list of statements).
func Interpret(statements []stmt.Stmt, errorReport *errorreport.ErrorReport) {
	defer func() {
		if e := recover(); e != nil {
			// This will intentionally re-panic if it's not a runtime error.
			runtimeError := e.(runtimeError)
			errorReport.ReportRuntimeError(runtimeError.token.Line, runtimeError.message)
		}
	}()
	i := interpreter{}
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i interpreter) VisitBinary(binary *expr.Binary) interface{} {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)

	switch binary.Operator.TokenType {
	case toks.Star:
		checkNumberOperands(binary.Operator, left, right)
		return left.(float64) * right.(float64)
	case toks.Slash:
		checkNumberOperands(binary.Operator, left, right)
		return left.(float64) / right.(float64)
	case toks.Minus:
		checkNumberOperands(binary.Operator, left, right)
		return left.(float64) - right.(float64)
	case toks.Plus:
		{
			l, leftOk := left.(float64)
			r, rightOk := right.(float64)

			if leftOk && rightOk {
				return l + r
			}
		}

		{
			l, leftOk := left.(string)
			r, rightOk := right.(string)

			if leftOk && rightOk {
				return l + r
			}
		}

		panic(runtimeError{token: binary.Operator, message: "Operands must be two numbers or two strings."})
	case toks.Greater:
		checkNumberOperands(binary.Operator, left, right)
		return left.(float64) > right.(float64)
	case toks.GreaterEqual:
		checkNumberOperands(binary.Operator, left, right)
		return left.(float64) >= right.(float64)
	case toks.Less:
		checkNumberOperands(binary.Operator, left, right)
		return left.(float64) < right.(float64)
	case toks.LessEqual:
		checkNumberOperands(binary.Operator, left, right)
		return left.(float64) <= right.(float64)
	case toks.EqualEqual:
		return isEqual(left, right)
	case toks.BangEqual:
		return !isEqual(left, right)
	}

	// Unreachable.
	return nil
}

func (i interpreter) VisitGrouping(grouping *expr.Grouping) interface{} {
	return i.evaluate(grouping.Expression)
}

func (i interpreter) VisitLiteral(literal *expr.Literal) interface{} {
	return literal.Value
}

func (i interpreter) VisitUnary(unary *expr.Unary) interface{} {
	right := i.evaluate(unary.Right)

	if unary.Operator.TokenType == toks.Bang {
		return !isTruthy(right)
	} else if unary.Operator.TokenType == toks.Minus {
		checkNumberOperand(unary.Operator, right)
		return -right.(float64)
	}

	// Unreachable.
	return nil
}

func (i interpreter) VisitStatementPrint(p *stmt.Print) {
	val := i.evaluate(p.Expression)
	fmt.Println(stringify(val))
}

func (i interpreter) VisitStatementExpression(e *stmt.Expression) {
	i.evaluate(e.Expression)
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}

	return true
}

func isEqual(left interface{}, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	return left == right
}

func (i interpreter) evaluate(expression expr.Expr) interface{} {
	// str := tools.PrintAst(expression)
	// fmt.Println(str)

	return expression.Accept(i)
}

func (i interpreter) execute(statement stmt.Stmt) {
	statement.Accept(i)
}

func checkNumberOperand(operator toks.Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}

	panic(runtimeError{token: operator, message: "Operand must be a number."})
}

func checkNumberOperands(operator toks.Token, operand1 interface{}, operand2 interface{}) {
	_, ok1 := operand1.(float64)
	_, ok2 := operand1.(float64)

	if ok1 && ok2 {
		return
	}

	panic(runtimeError{token: operator, message: "Operands must be numbers."})
}

func stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", val)
}
