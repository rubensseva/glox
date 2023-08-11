package main

import "fmt"

type Environment struct {
	enclosing *Environment // https://craftinginterpreters.com/statements-and-state.html#nesting-and-shadowing
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (any, error) {
	val, ok := e.values[name.lexeme]
	if ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	// return nil, RuntimeError{
	// 	token: name,
	// 	msg:   fmt.Sprintf("could not find token: %v in env: %v", name.lexeme, e.values),
	// }
	panic(RuntimeError{
		token: name,
		msg:   fmt.Sprintf("could not find token: %v in env: %v", name.lexeme, e.values),
	})
}

func (e *Environment) assign(name Token, value any) error {
	_, ok := e.values[name.lexeme]
	if ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return RuntimeError{
		token: name,
		msg:   fmt.Sprintf("Tried to assign undefined variable: Undefined variable '%s'.", name.lexeme),
	}
}
