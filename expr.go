package main

type Expr interface {
	Eval() Expr
}

// ASSIGN
type Assign struct {
	name  Token
	value Expr
}

func (b Assign) Eval() Expr {
	panic("assign not implemented yet")
}

// BINARY
type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) Eval() Expr {
	panic("binary eval not implemented yet")
}

// CALL
type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (b Call) Eval() Expr {
	panic("call eval not implemented yet")
}

// GET
type Get struct {
	object Expr
	name   Token
}

func (b Get) Eval() Expr {
	panic("get eval not implemented yet")
}

// GROUPING (
type Grouping struct {
	expression Expr
}

func (b Grouping) Eval() Expr {
	panic("grouping eval not implemented yet")
}

// LITERAL
type Literal struct {
	value any
}

func (b Literal) Eval() Expr {
	panic("literal eval not implemented yet")
}

// LOGICAL
type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Logical) Eval() Expr {
	panic("logical eval not implemented yet")
}

// UNARY
type Unary struct {
	operator Token
	right    Expr
}

func (b Unary) Eval() Expr {
	panic("logical eval not implemented yet")
}

// VARIABLE
type Variable struct {
	name Token
}

func (b Variable) Eval() Expr {
	panic("logical eval not implemented yet")
}
