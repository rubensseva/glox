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

type BlockStmt struct {
	statements []Stmt
}

func (s BlockStmt) IsStmt() {
	panic("shouldn't be called")
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (s IfStmt) IsStmt() {
	panic("shouldn't be called")
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (s WhileStmt) IsStmt() {
	panic("shouldn't be called")
}
