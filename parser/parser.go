/*
Package parser is used to convert a list of tokens to an abstract syntax tree using the following rules:

expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "false" | "true" | "nil"
			   | "(" expression ")" ;   */
package parser

import (
	"fmt"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/expr"
	"github.com/maleksiuk/golox/stmt"
	"github.com/maleksiuk/golox/toks"
)

type parser struct {
	current     int
	tokens      []toks.Token
	errorReport *errorreport.ErrorReport
}

type parseError struct {
	token   toks.Token
	message string
}

func (err *parseError) Error() string {
	return err.message
}

func newParseError(token toks.Token, message string) error {
	return &parseError{token: token, message: message}
}

// Parse converts a list of tokens to a list of statements.
func Parse(tokens []toks.Token, errorReport *errorreport.ErrorReport) []stmt.Stmt {
	p := parser{current: 0, tokens: tokens, errorReport: errorReport}

	var statements []stmt.Stmt

	for !p.isAtEnd() {
		statement, err := p.statement()

		// TODO: Eventually we will print out errors as they happen but for now they all bubble up and we'll print them here.
		if err != nil {
			if parseErr, ok := err.(*parseError); ok {
				p.printError(parseErr.token, parseErr.message)
			}

			return nil
		}

		statements = append(statements, statement)
	}

	return statements
}

func (p *parser) printError(token toks.Token, message string) {
	if token.TokenType == toks.EOF {
		p.errorReport.Report(token.Line, " at end", message)
	} else {
		p.errorReport.Report(token.Line, fmt.Sprintf(" at '%v'", token.Lexeme), message)
	}
}

func (p *parser) statement() (stmt.Stmt, error) {
	if p.match(toks.Print) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *parser) printStatement() (stmt.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(toks.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}

	return &stmt.Print{Expression: val}, nil
}

func (p *parser) expressionStatement() (stmt.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(toks.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}

	return &stmt.Expression{Expression: val}, nil
}

func (p *parser) expression() (expr.Expr, error) {
	return p.equality()
}

func (p *parser) previous() toks.Token {
	return p.tokens[p.current-1]
}

func (p *parser) equality() (expr.Expr, error) {
	expression, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(toks.BangEqual, toks.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (p *parser) comparison() (expr.Expr, error) {
	expression, err := p.addition()
	if err != nil {
		return nil, err
	}

	for p.match(toks.Greater, toks.GreaterEqual, toks.Less, toks.LessEqual) {
		operator := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}

		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (p *parser) addition() (expr.Expr, error) {
	expression, err := p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(toks.Minus, toks.Plus) {
		operator := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (p *parser) multiplication() (expr.Expr, error) {
	expression, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(toks.Star, toks.Slash) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (p *parser) unary() (expr.Expr, error) {
	if p.match(toks.Bang, toks.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &expr.Unary{Operator: operator, Right: right}, nil
	}

	return p.primary()
}

func (p *parser) primary() (expr.Expr, error) {
	if p.match(toks.Number, toks.String) {
		return &expr.Literal{Value: p.previous().Literal}, nil
	}

	if p.match(toks.False) {
		return &expr.Literal{Value: false}, nil
	}

	if p.match(toks.True) {
		return &expr.Literal{Value: true}, nil
	}

	if p.match(toks.Nil) {
		return &expr.Literal{Value: nil}, nil
	}

	if p.match(toks.LeftParen) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(toks.RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return &expr.Grouping{Expression: expression}, nil
	}

	return nil, newParseError(p.peek(), "expect expression")
}

func (p *parser) consume(tokenType toks.TokenType, errorMessage string) (toks.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return toks.Token{}, newParseError(p.peek(), errorMessage)
}

func (p *parser) advance() toks.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *parser) match(tokenTypes ...toks.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.peek().TokenType == tokenType {
			p.advance()
			return true
		}
	}

	return false
}

func (p *parser) check(tokenType toks.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == tokenType
}

func (p *parser) peek() toks.Token {
	return p.tokens[p.current]
}

func (p *parser) isAtEnd() bool {
	return p.peek().TokenType == toks.EOF
}
