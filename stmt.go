package main

type Stmt interface {
	IsStmt()
}

type PrintStmt struct {
	expression Expr
}

func (s PrintStmt) IsStmt() {
	panic("shouldn't be called")
}

type ExpressionStmt struct {
	expression Expr
}

func (s ExpressionStmt) IsStmt() {
	panic("shouldn't be called")
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (s VarStmt) IsStmt() {
	panic("shouldn't be called")
}
