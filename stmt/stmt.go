package stmt

import (
	"github.com/maleksiuk/golox/expr"
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

type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(visitor Visitor) {
	visitor.VisitStatementPrint(p)
}

type Visitor interface {
	VisitStatementExpression(expression *Expression)
	VisitStatementPrint(p *Print)
}
