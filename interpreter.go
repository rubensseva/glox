package main

import (
	"fmt"
)

func interpret(expression Expr) {
	var value any = evaluate(expression)
	fmt.Println(value)
}

func evaluate(expr Expr) any {
	switch t := expr.(type) {
	case Binary:
		return evalBinary(t)
	case Grouping:
		return evalGrouping(t)
	case Literal:
		return evalLiteral(t)
	case Unary:
		return evalUnary(t)
	default:
		panic(fmt.Sprintf("unknown type %T: %v", expr, t))
	}
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if res, ok := object.(bool); ok {
		return res
	}
	return true
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	// TODO I think this should work https://go.dev/play/p/0T79bqDjO8_B
	// But not 100% sure it works properly for all cases
	return a == b
}

func evalBinary(expr Binary) any {
	var left any = evaluate(expr.left)
	var right any = evaluate(expr.right)

	switch expr.operator.tokenType {
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case MINUS:
		return left.(float64) - right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case PLUS:
		// Pluss is a bit special because it works for
		// numbers and strings
		{
			l, lok := left.(float64)
			r, rok := right.(float64)
			if lok && rok {
				return l + r
			}
		}
		{
			l, lok := left.(string)
			r, rok := right.(string)
			if lok && rok {
				return l + r
			}
		}
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	}

	// Unreachable
	return nil
}

func evalGrouping(expr Grouping) any {
	return evaluate(expr.expression)
}

func evalLiteral(expr Literal) any {
	return expr.value
}

func evalUnary(expr Unary) any {
	var right any = evaluate(expr.right)

	switch expr.operator.tokenType {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		return -(right.(float64))
	}

	// Unreachable
	return nil
}
