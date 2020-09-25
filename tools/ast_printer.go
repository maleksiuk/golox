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

func (printer astPrinter) VisitLogical(logical *expr.Logical) interface{} {
	return printer.parenthesize(logical.Operator.Lexeme, logical.Left, logical.Right)
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

func (printer astPrinter) VisitVariable(v *expr.Variable) interface{} {
	return v.Name.Lexeme
}

func (printer astPrinter) VisitAssign(assign *expr.Assign) interface{} {
	return printer.parenthesize("=", assign.Name.Lexeme, assign.Value)
}

func (printer astPrinter) VisitCall(call *expr.Call) interface{} {
	var args strings.Builder

	for idx, ele := range call.Arguments {
		args.WriteString(ele.Accept(printer).(string))
		if idx != len(call.Arguments)-1 {
			args.WriteString(",")
		}
	}

	if args.Len() > 0 {
		return printer.parenthesize("call", call.Callee, args.String())
	} else {
		return printer.parenthesize("call", call.Callee)
	}
}

func (printer astPrinter) parenthesize(name string, parts ...interface{}) string {
	var str strings.Builder

	str.WriteString("(")
	str.WriteString(name)

	for _, part := range parts {
		str.WriteString(" ")
		switch p := part.(type) {
		case expr.Expr:
			str.WriteString(p.Accept(printer).(string))
		case string:
			str.WriteString(p)
		case fmt.Stringer:
			str.WriteString(p.String())
		}
	}
	str.WriteString(")")

	return str.String()
}
