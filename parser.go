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

func (p *Parser) whileStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.");
	condition := p.expression();
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")

	body, err := p.statement()
	if err != nil {
		return nil, fmt.Errorf("getting statement for body: %w", err)
	}
	return WhileStmt{
		condition: condition,
		body:      body,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}

	// https://craftinginterpreters.com/statements-and-state.html#block-syntax-and-semantics
	if p.match(LEFT_BRACE) {
		return BlockStmt{
			statements: p.block(),
		}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	var zero Stmt = nil
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	var initializerIsSet bool

	if p.match(SEMICOLON) {
		initializerIsSet = false
		// initializer = nil
	} else if p.match(VAR) {
		tmp, err := p.varDeclaration()
		if err != nil {
			return zero, fmt.Errorf("handling initializer var declaration: %w", err)
		}
		initializer = tmp
		initializerIsSet = true
	} else {
		tmp, err := p.expressionStatement()
		if err != nil {
			return zero, fmt.Errorf("handling initializer expr stmt: %w", err)
		}
		initializer = tmp
		initializerIsSet = true
	}

	var condition Expr
	var conditionIsSet bool
	if !p.check(SEMICOLON) {
		condition = p.expression()
		conditionIsSet = true
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.");

	var increment Expr
	var incrementIsSet bool
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
		incrementIsSet = true
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body, err := p.statement()
	if err != nil {
		return zero, fmt.Errorf("getting statement in for loop desugaring: %w", err)
	}

	if incrementIsSet {
		body = BlockStmt{
			statements: []Stmt{
				body,
				ExpressionStmt{
					expression: increment,
				},
			},
		}
	}

	// Note reverse check
	if !conditionIsSet {
		condition = Literal{
			value: true,
		}
	}

	body = WhileStmt{
		condition: condition,
		body:      body,
	}


	if initializerIsSet {
		body = BlockStmt{
			statements: []Stmt{initializer, body},
		}
	}

	return body, nil

}

func (p *Parser) ifStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after if condition.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, fmt.Errorf("getting statement for then branch: %w", err)
	}
	// TODO: Careful! Using nil interface here.. does it work?
	var elseBranch Stmt = nil
	if (p.match(ELSE)) {
		tmp, err := p.statement()
		if err != nil {
			return nil, fmt.Errorf("getting statement for else branch: %w", err)
		}
		elseBranch = tmp
	}

	return IfStmt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}, nil
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

func (p *Parser) block() []Stmt {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(
			statements,
			p.declaration(),
		)
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) assignment() Expr {
	var expr Expr = p.or()

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

func (p *Parser) or() Expr {
	var expr Expr = p.and()

	for p.match(OR) {
		var operator Token = p.previous()
		var right Expr = p.and()
		expr = Logical{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) and() Expr {
	var expr Expr = p.equality()

	for p.match(AND) {
		var operator Token = p.previous()
		var right Expr = p.equality()
		expr = Logical{
			left:     expr,
			operator: operator,
			right:    right,
		}
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
