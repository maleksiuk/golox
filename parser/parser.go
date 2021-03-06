/*
Package parser is used to convert a list of tokens to an abstract syntax tree using the following rules:

expression     → assignment ;
assignment     → IDENTIFIER "=" assignment
			   | logic_or ;
logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
			   | call ;
call           → primary ( "(" arguments? ")" )* ;
arguments      → expression ( "," expression )* ;
primary        → NUMBER | STRING | "false" | "true" | "nil"
			   | "(" expression ")"
			   | IDENTIFIER ;

program     → declaration* EOF ;
declaration → funDecl
            | varDecl
			| statement ;
varDecl     → "var" IDENTIFIER ( "=" expression )? ";" ;
funDecl        → "fun" function ;
function       → IDENTIFIER "(" parameters? ")" block ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
statement → exprStmt
          | ifStmt
          | printStmt
		  | whileStmt
		  | forStmt
          | block ;

exprStmt  → expression ";" ;
printStmt → "print" expression ";" ;
whileStmt → "while" "(" expression ")" statement ;
forStmt   → "for" "(" ( varDecl | exprStmt | ";" )
                      expression? ";"
                      expression? ")" statement ;
ifStmt    → "if" "(" expression ")" statement ( "else" statement )? ;
*/
package parser

// TODO: panic instead of passing errors all the way up the chain

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
		statement, err := p.declaration()

		// TODO: Eventually all errors will be printed out as they happen and this block will be deleted.
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
		p.errorReport.Report(token.Line, "at end", message)
	} else {
		p.errorReport.Report(token.Line, fmt.Sprintf("at '%v'", token.Lexeme), message)
	}
}

func (p *parser) declaration() (stmt.Stmt, error) {
	defer func() {
		if e := recover(); e != nil {
			_, ok := e.(*parseError)
			if !ok {
				panic(e)
			}

			p.synchronize()
			return
		}
	}()

	if p.match(toks.Fun) {
		return p.funDeclaration()
	}

	if p.match(toks.Var) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *parser) varDeclaration() (stmt.Stmt, error) {
	nameToken := p.consume(toks.Identifier, "Expect variable name.")

	var expr expr.Expr
	var err error
	if p.match(toks.Equal) {
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(toks.Semicolon, "Expect ';' after variable declaration.")

	return &stmt.Var{Name: nameToken, Initializer: expr}, nil
}

func (p *parser) funDeclaration() (stmt.Stmt, error) {
	nameToken := p.consume(toks.Identifier, "Expect function name.")

	p.consume(toks.LeftParen, "Expect '(' after function name")

	parameters := make([]toks.Token, 0, 10)
	matchedComma := true

	if !p.check(toks.RightParen) {
		for matchedComma {
			if len(parameters) >= 255 {
				// TODO: This is one that should just be printed out here without stopping parsing.
				return nil, newParseError(p.peek(), "Cannot have more than 255 parameters.")
			}

			identifierToken := p.consume(toks.Identifier, "Expect parameter name.")
			parameters = append(parameters, identifierToken)

			matchedComma = p.match(toks.Comma)
		}
	}

	p.consume(toks.RightParen, "Expect ')' after parameters")

	p.consume(toks.LeftBrace, "Expect '{' before function body.")
	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &stmt.Function{Name: nameToken, Params: parameters, Body: body}, nil
}

func (p *parser) statement() (stmt.Stmt, error) {
	if p.match(toks.Print) {
		return p.printStatement()
	}

	if p.match(toks.While) {
		return p.whileStatement()
	}

	if p.match(toks.For) {
		return p.forStatement()
	}

	if p.match(toks.If) {
		return p.conditionalStatement()
	}

	if p.match(toks.LeftBrace) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &stmt.Block{Statements: statements}, nil
	}

	return p.expressionStatement()
}

func (p *parser) block() ([]stmt.Stmt, error) {
	var statements = make([]stmt.Stmt, 0, 10)

	for !p.check(toks.RightBrace) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	p.consume(toks.RightBrace, "Expect '}' after block.")

	return statements, nil
}

func (p *parser) conditionalStatement() (stmt.Stmt, error) {
	p.consume(toks.LeftParen, "Expect '(' after 'if'.")

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(toks.RightParen, "Expect ')' after 'if' condition.")

	thenStatement, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseStatement stmt.Stmt
	if p.match(toks.Else) {
		elseStatement, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &stmt.Conditional{Condition: condition, ThenStatement: thenStatement, ElseStatement: elseStatement}, nil
}

func (p *parser) printStatement() (stmt.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(toks.Semicolon, "Expect ';' after value")

	return &stmt.Print{Expression: val}, nil
}

func (p *parser) whileStatement() (stmt.Stmt, error) {
	p.consume(toks.LeftParen, "Expect '(' after while")

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(toks.RightParen, "Expect ')' after while")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &stmt.While{Condition: condition, Body: body}, nil
}

/*
forStmt   → "for" "(" ( varDecl | exprStmt | ";" )
                      expression? ";"
					  expression? ")" statement ;
*/
func (p *parser) forStatement() (stmt.Stmt, error) {
	var err error

	p.consume(toks.LeftParen, "Expect '(' after for")

	var initializer stmt.Stmt
	var condition expr.Expr
	var increment expr.Expr

	if p.match(toks.Semicolon) {
		initializer = nil
	} else if p.match(toks.Var) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	if p.match(toks.Semicolon) {
		condition = nil
	} else {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(toks.Semicolon, "Expect ';' after for condition")
	}

	if p.match(toks.RightParen) {
		increment = nil
	} else {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(toks.RightParen, "Expect ')' at end of for")
	}

	userSpecifiedBody, err := p.statement()
	if err != nil {
		return nil, err
	}

	// the body is the user-specified body followed by the (optional) increment
	var bodyStatements = make([]stmt.Stmt, 0, 2)
	bodyStatements = append(bodyStatements, userSpecifiedBody)
	if increment != nil {
		bodyStatements = append(bodyStatements, &stmt.Expression{Expression: increment})
	}
	bodyBlock := &stmt.Block{Statements: bodyStatements}

	if condition == nil {
		condition = &expr.Literal{Value: true}
	}

	while := &stmt.While{Condition: condition, Body: bodyBlock}

	// the final result is the (optional) initializer followed by the while loop
	var statements = make([]stmt.Stmt, 0, 2)
	if initializer != nil {
		statements = append(statements, initializer)
	}
	statements = append(statements, while)
	block := &stmt.Block{Statements: statements}

	return block, nil
}

func (p *parser) expressionStatement() (stmt.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(toks.Semicolon, "Expect ';' after value")

	return &stmt.Expression{Expression: val}, nil
}

func (p *parser) expression() (expr.Expr, error) {
	return p.assignment()
}

func (p *parser) previous() toks.Token {
	return p.tokens[p.current-1]
}

func (p *parser) logicAnd() (expr.Expr, error) {
	expression, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(toks.And) {
		andToken := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expression = &expr.Logical{Left: expression, Operator: andToken, Right: right}
	}

	return expression, nil
}

func (p *parser) logicOr() (expr.Expr, error) {
	expression, err := p.logicAnd()
	if err != nil {
		return nil, err
	}

	for p.match(toks.Or) {
		orToken := p.previous()
		right, err := p.logicAnd()
		if err != nil {
			return nil, err
		}

		expression = &expr.Logical{Left: expression, Operator: orToken, Right: right}
	}

	return expression, nil
}

func (p *parser) assignment() (expr.Expr, error) {
	expression, err := p.logicOr()
	if err != nil {
		return nil, err
	}

	if p.match(toks.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expression.(*expr.Variable); ok {
			return &expr.Assign{Name: variable.Name, Value: value}, nil
		}

		p.handleError(equals, "Invalid assignment target")
	}

	return expression, nil
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

	return p.call()
}

func (p *parser) call() (expr.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	// the book says doing true/break will be better later on
	for true {
		if p.match(toks.LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *parser) finishCall(callee expr.Expr) (expr.Expr, error) {
	args := make([]expr.Expr, 0, 5)

	if !p.check(toks.RightParen) {
		matchedComma := true

		for matchedComma {
			// the book implementation also has >= here even though the error message implies >
			if len(args) >= 255 {
				p.handleError(p.peek(), "Cannot have more than 255 arguments.")
			}

			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			matchedComma = p.match(toks.Comma)
		}
	}

	rightParenToken := p.consume(toks.RightParen, "Expected ')' to finish function call")

	return &expr.Call{Callee: callee, Paren: rightParenToken, Arguments: args}, nil
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

		p.consume(toks.RightParen, "Expect ')' after expression.")

		return &expr.Grouping{Expression: expression}, nil
	}

	if p.match(toks.Identifier) {
		return &expr.Variable{Name: p.previous()}, nil
	}

	return nil, newParseError(p.peek(), "expect expression")
}

func (p *parser) consume(tokenType toks.TokenType, errorMessage string) toks.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	err := p.handleError(p.peek(), errorMessage)
	panic(err)
}

// correspond's to jlox's error() function
func (p *parser) handleError(token toks.Token, message string) error {
	p.printError(token, message)
	return newParseError(token, message)
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

func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == toks.Semicolon {
			return
		}

		switch p.peek().TokenType {
		case toks.Class, toks.Fun, toks.Var, toks.For, toks.If, toks.While, toks.Print, toks.Return:
			return
		}

		p.advance()
	}
}
