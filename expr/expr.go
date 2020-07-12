package expr

import (
	"github.com/maleksiuk/golox/toks"
)

type Expr interface {
	Accept(visitor Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator toks.Token
	Right    Expr
}

func (binary *Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinary(binary)
}

type Grouping struct {
	Expression Expr
}

func (grouping *Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGrouping(grouping)
}

type Literal struct {
	Value interface{}
}

func (literal *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(literal)
}

type Unary struct {
	Operator toks.Token
	Right    Expr
}

func (unary *Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnary(unary)
}

type Visitor interface {
	VisitBinary(binary *Binary) interface{}
	VisitGrouping(grouping *Grouping) interface{}
	VisitLiteral(literal *Literal) interface{}
	VisitUnary(unary *Unary) interface{}
}