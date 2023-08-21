package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration FunctionStmt
	closure     *Environment
}

// Type check, just to be safe
var _ LoxCallable = &LoxFunction{}

func NewLoxFunction(declaration FunctionStmt, closure *Environment) *LoxFunction {
	return &LoxFunction{
		closure:     closure,
		declaration: declaration,
	}
}

// Using named return values here so we can modify the returned value in deferred function
func (l *LoxFunction) Call(interpreter *Interpreter, arguments []any) (result any, err error) {
	environment := NewEnvironment(l.closure)
	for i := 0; i < len(l.declaration.params); i++ {
		environment.define(
			l.declaration.params[i].lexeme,
			arguments[i],
		)
	}

	defer func() {
		if val := recover(); val != nil {
			v, ok := val.(ReturnHack)
			if !ok {
				panic(val)
			}
			// HACK: Modify the return value
			// See https://yourbasic.org/golang/defer/
			result = v.value
			err = nil
		}
	}()
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
