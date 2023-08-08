package main

import (
	"fmt"
	"time"
)

type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(Interpreter Interpreter, arguments []any) (any, error) {
	return time.Now().Unix(), nil
}

func (c *Clock) String() string {
	return fmt.Sprintf("<native fn Clock()>")
}
