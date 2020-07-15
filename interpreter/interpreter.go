package interpreter

import (
	"fmt"

	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/toks"
)

type interpreter struct {
}

// Interpret evaluates an expression and prints out the result.
func Interpret(expression expr.Expr) {
	i := interpreter{}
	fmt.Println(i.evaluate(expression))
}

func (i interpreter) VisitBinary(binary *expr.Binary) interface{} {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)

	switch binary.Operator.TokenType {
	case toks.Star:
		return left.(float64) * right.(float64)
	case toks.Slash:
		return left.(float64) / right.(float64)
	case toks.Minus:
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
	case toks.Greater:
		return left.(float64) > right.(float64)
	case toks.GreaterEqual:
		return left.(float64) >= right.(float64)
	case toks.Less:
		return left.(float64) < right.(float64)
	case toks.LessEqual:
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
		return -right.(float64)
	}

	// Unreachable.
	return nil
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
	return expression.Accept(i)
}
