package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration FunctionStmt
}

// Type check, just to be safe
var _ LoxCallable = &LoxFunction{}

func NewLoxFunction(declaration FunctionStmt) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
	}
}

func (l *LoxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironment(interpreter.globals)
	for i := 0; i < len(l.declaration.params); i++ {
		environment.define(
			l.declaration.params[i].lexeme,
			arguments[i],
		)
	}
	if err := interpreter.executeBlock(l.declaration.body, environment); err != nil {
		return nil, fmt.Errorf("executing block: %w", err)
	}
	return nil, nil
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.params)
}

func (l *LoxFunction) String() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}
