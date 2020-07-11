// Package tools implements types and functions that help the developer debug the program.
package tools

import (
	"fmt"
	"strings"

	"github.com/maleksiuk/golox/expr"
)

// PrintAst returns a string representation of the AST represented by expression.
func PrintAst(expression expr.Expr) string {
	printer := astPrinter{}
	return printer.Print(expression)
}

type astPrinter struct {
}

func (printer astPrinter) Print(expression expr.Expr) string {
	return expression.Accept(printer).(string)
}

func (printer astPrinter) VisitBinary(binary *expr.Binary) interface{} {
	return printer.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (printer astPrinter) VisitGrouping(grouping *expr.Grouping) interface{} {
	return printer.parenthesize("group", grouping.Expression)
}

func (printer astPrinter) VisitLiteral(literal *expr.Literal) interface{} {
	return fmt.Sprintf("%v", literal.Value)
}

func (printer astPrinter) VisitUnary(unary *expr.Unary) interface{} {
	return printer.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (printer astPrinter) parenthesize(name string, expressions ...expr.Expr) string {
	var str strings.Builder

	str.WriteString("(")
	str.WriteString(name)

	for _, expression := range expressions {
		str.WriteString(" ")
		str.WriteString(expression.Accept(printer).(string))
	}
	str.WriteString(")")

	return str.String()
}
