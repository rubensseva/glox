package main

import (
	"fmt"
)

type RuntimeError struct {
	token Token
	msg   string
	err   error
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf(
		"line: %v, lexeme: %v, tokentype: %v, literal: %v, msg: %v",
		r.token.line,
		r.token.lexeme,
		r.token.tokenType,
		r.token.literal,
		r.msg,
	)
}

func (r RuntimeError) Unwrap() error {
	return r.err
}
