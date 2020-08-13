package stmt

import (
	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/toks"
)

type Stmt interface {
	Accept(visitor Visitor)
}

type Expression struct {
	Expression expr.Expr
}

func (expression *Expression) Accept(visitor Visitor) {
	visitor.VisitStatementExpression(expression)
}

type Block struct {
	Statements []Stmt
}

func (block *Block) Accept(visitor Visitor) {
	visitor.VisitBlock(block)
}

type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(visitor Visitor) {
	visitor.VisitStatementPrint(p)
}

type Var struct {
	Name        toks.Token
	Initializer expr.Expr
}

func (v *Var) Accept(visitor Visitor) {
	visitor.VisitStatementVar(v)
}

type Visitor interface {
	VisitStatementExpression(expression *Expression)
	VisitStatementPrint(p *Print)
	VisitStatementVar(v *Var)
	VisitBlock(block *Block)
}
