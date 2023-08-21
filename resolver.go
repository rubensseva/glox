package main

import (
	"fmt"
)

type Resolver struct {
	scopes Stack[map[string]bool]
	interpreter *Interpreter
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
	}
}

// The book uses function overloading to have a single
// resolve statement. We dont have that, so we need both
// resolveStatements and resolveStatement.
// See chapter 11.3.1
func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}
func (r *Resolver) resolveStatement(statement Stmt) {
	switch t := statement.(type) {
	case PrintStmt:
		r.visitPrintStmt(t)
	case VarStmt:
		r.visitVarStmt(t)
	case ExpressionStmt:
		r.visitExpressionStmt(t)
	case BlockStmt:
		r.visitBlockStmt(t)
	case IfStmt:
		r.visitIfStmt(t)
	case WhileStmt:
		r.visitWhileStmt(t)
	case FunctionStmt:
		r.visitFunctionStmt(t)
	case ReturnStmt:
		r.visitReturnStmt(t)
	default:
		panic(fmt.Errorf("couldnt find type for stmt: %[1]v %[1]T", statement))
	}
}
func (r *Resolver) resolveExpr(expr Expr) {
	switch t := expr.(type) {
	case Binary:
		r.visitBinaryExpr(t)
	case Grouping:
		r.visitGroupingExpr(t)
	case Literal:
		r.visitLiteralExpr(t)
	case Unary:
		r.visitUnaryExpr(t)
	case Variable:
		r.visitVariableExpr(t)
	case Logical:
		r.visitLogicalExpr(t)
	case Assign:
		r.visitAssignExpr(t)
	case Call:
		r.visitCallExpr(t)
	default:
		panic(fmt.Errorf("couldnt find type for expr: %[1]v %[1]T", expr))
	}
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	v := r.scopes.Pop()
	fmt.Printf("stack popped value: %v\n", v)
}

func (r *Resolver) declare(name Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek()
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if r.scopes.IsEmpty() {
		return
	}
	r.scopes.Peek()[name.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		_, ok := r.scopes.Get(i)[name.lexeme]
		if ok {
			// TODO: Uncomment
			// r.interpreter.resolve(expr, r.scopes.Len() - 1 - i)
			return
		}
	}
}

func (r *Resolver) visitBlockStmt(stmt BlockStmt) {
	r.beginScope()
	r.resolveStatements(stmt.statements)
	r.endScope()
}

func (r *Resolver) visitExpressionStmt(stmt ExpressionStmt) {
	r.resolveExpr(stmt.expression)
}

func (r *Resolver) visitFunctionStmt(stmt FunctionStmt) {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt)
}

func (r *Resolver) visitIfStmt(stmt IfStmt) {
	r.resolveExpr(stmt.condition)
	r.resolveStatement(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStatement(stmt.elseBranch)
	}
}

func (r *Resolver) visitPrintStmt(stmt PrintStmt) {
	r.resolveExpr(stmt.expression)
}

func (r *Resolver) visitReturnStmt(stmt ReturnStmt) {
	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
}

func (r *Resolver) resolveFunction(function FunctionStmt) {
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.body)
	r.endScope()
}

func (r *Resolver) visitVarStmt(stmt VarStmt) {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpr(stmt.initializer)
	}
	r.define(stmt.name)
}

func (r *Resolver) visitWhileStmt(stmt WhileStmt) {
	r.resolveExpr(stmt.condition)
	r.resolveStatement(stmt.body)
}

func (r *Resolver) visitAssignExpr(expr Assign) {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
}

func (r *Resolver) visitBinaryExpr(expr Binary) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
}

func (r *Resolver) visitCallExpr(expr Call) {
	r.resolveExpr(expr.callee)

	for _, argument := range expr.arguments {
		r.resolveExpr(argument)
	}
}

func (r *Resolver) visitGroupingExpr(expr Grouping) {
	r.resolveExpr(expr.expression)
}

func (r *Resolver) visitLiteralExpr(expr Literal) {
	// nothing to do
}

func (r *Resolver) visitLogicalExpr(expr Logical) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
}

func (r *Resolver) visitUnaryExpr(expr Unary) {
	r.resolveExpr(expr.right)
}

func (r *Resolver) visitVariableExpr(expr Variable) {
	if !r.scopes.IsEmpty() &&
		r.scopes.Peek()[expr.name.lexeme] == false {
		// loxtokenerror(expr.name,
		// 	"Cant't read local variable in its own initializer.");
		panic(fmt.Errorf("Cant't read local variable in its own initializer. expr name: %v", expr.name));
	}
	r.resolveLocal(expr, expr.name)
}
