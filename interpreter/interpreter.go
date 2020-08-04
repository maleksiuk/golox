package interpreter

import (
	"fmt"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/stmt"
	"github.com/maleksiuk/golox/toks"
)

type environment struct {
	variables map[string]interface{}
}

func (e *environment) define(name string, val interface{}) {
	e.variables[name] = val
}

func (e *environment) assign(name toks.Token, val interface{}) {
	nameValue := name.Lexeme
	_, ok := e.variables[nameValue]
	if ok {
		e.variables[nameValue] = val
	} else {
		message := fmt.Sprintf("Undefined variable '%v'.", nameValue)
		panic(runtimeError{token: name, message: message})
	}
}

func (e *environment) get(name toks.Token) interface{} {
	val, ok := e.variables[name.Lexeme]
	if !ok {
		message := fmt.Sprintf("Undefined variable '%v'.", name.Lexeme)
		panic(runtimeError{token: name, message: message})
	}

	return val
}

type Interpreter struct {
	env *environment
}

type runtimeError struct {
	token   toks.Token
	message string
}

func NewInterpreter() Interpreter {
	env := environment{variables: make(map[string]interface{})}
	return Interpreter{env: &env}
}

// Interpret evaluates a program (list of statements).
func (i Interpreter) Interpret(statements []stmt.Stmt, errorReport *errorreport.ErrorReport) {
	defer func() {
		if e := recover(); e != nil {
			// This will intentionally re-panic if it's not a runtime error.
			runtimeError := e.(runtimeError)
			errorReport.ReportRuntimeError(runtimeError.token.Line, runtimeError.message)
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i Interpreter) VisitBinary(binary *expr.Binary) interface{} {
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

func (i Interpreter) VisitGrouping(grouping *expr.Grouping) interface{} {
	return i.evaluate(grouping.Expression)
}

func (i Interpreter) VisitLiteral(literal *expr.Literal) interface{} {
	return literal.Value
}

func (i Interpreter) VisitUnary(unary *expr.Unary) interface{} {
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

func (i Interpreter) VisitVariable(v *expr.Variable) interface{} {
	return i.env.get(v.Name)
}

func (i Interpreter) VisitAssign(assign *expr.Assign) interface{} {
	value := i.evaluate(assign.Value)
	i.env.assign(assign.Name, value)

	return value
}

func (i Interpreter) VisitStatementPrint(p *stmt.Print) {
	val := i.evaluate(p.Expression)
	fmt.Println(stringify(val))
}

func (i Interpreter) VisitStatementExpression(e *stmt.Expression) {
	i.evaluate(e.Expression)
}

func (i Interpreter) VisitStatementVar(v *stmt.Var) {
	var val interface{} = nil
	if v.Initializer != nil {
		val = i.evaluate(v.Initializer)
	}

	i.env.define(v.Name.Lexeme, val)
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

func (i Interpreter) evaluate(expression expr.Expr) interface{} {
	// str := tools.PrintAst(expression)
	// fmt.Println(str)

	return expression.Accept(i)
}

func (i Interpreter) execute(statement stmt.Stmt) {
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
