package main

import (
	"fmt"
	"strings"
)

func ASTPrint(expr Expr) string {
	switch t := expr.(type) {
	case Binary:
		return printBinary(t)
	case Grouping:
		return printGrouping(t)
	case Literal:
		return printLiteral(t)
	case Unary:
		return printUnary(t)
	default:
		panic(fmt.Sprintf("unknown type %T: %v", expr, t))
	}
}

func parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
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

func printGrouping(expr Grouping) string {
	return parenthesize("group", expr.expression)
}

func printLiteral(expr Literal) string {
	if expr.value == nil {
		return "nil"
	}
	// TODO: Better way to do it than through fmt?
	return fmt.Sprintf("%v", expr.value)
}

func printUnary(expr Unary) string {
	return parenthesize(expr.operator.lexeme, expr.right)
}

//   @Override
//   public String visitGroupingExpr(Expr.Grouping expr) {
//     return parenthesize("group", expr.expression);
//   }

//   @Override
//   public String visitLiteralExpr(Expr.Literal expr) {
//     if (expr.value == null) return "nil";
//     return expr.value.toString();
//   }

//   @Override
//   public String visitUnaryExpr(Expr.Unary expr) {
//     return parenthesize(expr.operator.lexeme, expr.right);
//   }
// }
