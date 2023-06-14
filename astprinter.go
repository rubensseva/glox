package main

import (
	"fmt"
	"strings"
)

func ASTPrint(expr Expr) string {
	switch t := expr.(type) {
	case Binary:
		return printBinary(t)
	default:
		panic(fmt.Sprintf("unknown type %T: %v", expr, t))
	}
}

func parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(ASTPrint(expr))
	}
	builder.WriteString(")")
	return builder.String()
}

func printBinary(expr Binary) string {
	return parenthesize(expr.operator.lexeme, expr.left, expr.right)
}
