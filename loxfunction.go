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

// Using named return values here so we can modify the returned value in deferred function
func (l *LoxFunction) Call(interpreter *Interpreter, arguments []any) (result any, err error) {
	fmt.Printf("calling func %v with args %v\n", l.declaration.name, arguments)
	fmt.Printf("env: %+v\n", interpreter.globals)

	environment := NewEnvironment(interpreter.globals)
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
			fmt.Printf("got returnhack: %[1]T %[1]+v %[1]f\n", v.value)
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
