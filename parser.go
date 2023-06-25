package main

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p *Parser) parse() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt := p.declaration()
		statements = append(
			statements,
			stmt,
		)
	}
	return statements, nil
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) declaration() Stmt {
	// try
	var err error
	var res Stmt
	if p.match(VAR) {
		res, err = p.varDeclaration()
	} else {
		res, err = p.statement()
	}
	if err != nil {
		fmt.Printf("warning: encountered error but ignoring and synchronizing instead. The ignored error is: %v\n", err)
		p.synchronize()
		return nil
	}
	return res
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, fmt.Errorf("consuming identifier: %w", err)
	}
	var initializer Expr = nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, fmt.Errorf("consuming semicolon: %w", err)
	}
	return VarStmt{
		name:        name,
		initializer: initializer,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	var value Expr = p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, fmt.Errorf("consuming semicolon: %w", err)
	}
	return PrintStmt{
		expression: value,
	}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	var expr Expr = p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, fmt.Errorf("soncuming semicolon: %w", err)
	}
	return ExpressionStmt{
		expression: expr,
	}, nil
}

func (p *Parser) assignment() Expr {
	var expr Expr = p.equality()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		foo, ok := expr.(Variable)
		if ok {
			var name Token = foo.name
			return Assign{
				name:  name,
				value: value,
			}
		}

		loxtokenerror(equals, "Invalid assignment target.")
	}

	return expr
}

// match checks if the current token matches any of the types
// in the input arg(s). If it matches, consume the token
// and return true
func (p *Parser) match(types ...TokenType) bool {
	for _, tokentype := range types {
		if p.check(tokentype) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokentype TokenType, message string) (Token, error) {
	if p.check(tokentype) {
		return p.advance(), nil
	}

	err := p.error(p.peek(), message)
	return Token{}, err
}

func (p *Parser) error(token Token, message string) error {
	loxtokenerror(token, message)

	// TODO: Figure this out
	// https://craftinginterpreters.com/parsing-expressions.html#entering-panic-mode
	return fmt.Errorf("parse error: Token: %v, message: %v", token, message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

// check checks if the current token is of the token type
// in the argument
func (p *Parser) check(tokentype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokentype
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) equality() Expr {
	var expr Expr = p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		var operator Token = p.previous()
		var right Expr = p.comparison()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	var expr Expr = p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		var operator Token = p.previous()
		var right Expr = p.term()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	var expr Expr = p.factor()

	for p.match(MINUS, PLUS) {
		var operator Token = p.previous()
		var right Expr = p.factor()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	var expr Expr = p.unary()

	for p.match(SLASH, STAR) {
		var operator Token = p.previous()
		var right Expr = p.unary()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		var operator Token = p.previous()
		var right Expr = p.unary()
		return Unary{
			operator: operator,
			right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{
			value: false,
		}
	}
	if p.match(TRUE) {
		return Literal{
			value: true,
		}
	}
	if p.match(NIL) {
		return Literal{
			value: nil,
		}
	}

	if p.match(NUMBER, STRING) {
		return Literal{
			value: p.previous().literal,
		}
	}

	if p.match(IDENTIFIER) {
		return Variable{
			name: p.previous(),
		}
	}

	if p.match(LEFT_PAREN) {
		var expr Expr = p.expression()
		_, err := p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			panic(fmt.Errorf("trying to consume: %w", err))
		}
		return Grouping{
			expression: expr,
		}
	}

	// TODO: Should probaly return a custom error
	panic(p.error(p.peek(), "Expect expression."))
}
