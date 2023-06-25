package main

import "fmt"

type Environment struct {
	values map[string]any
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (any, error) {
	val, ok := e.values[name.lexeme]
	if !ok {
		return nil, RuntimeError{
			token: name,
			msg:   fmt.Sprintf("could not find token: %v in env: %v", name.lexeme, e.values),
		}
	}
	return val, nil
}

func (e *Environment) assign(name Token, value any) error {
	_, ok := e.values[name.lexeme]
	if ok {
		e.values[name.lexeme] = value
		return nil
	}

	return RuntimeError{
		token: name,
		msg:   fmt.Sprintf("Tried to assign undefined variable: Undefined variable '%s'.", name.lexeme),
	}
}
