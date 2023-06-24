package main

import (
	"errors"
	"fmt"
)

func interpret(expression Expr) {
	value, err := evaluate(expression)
	if err != nil {
		var trgt RuntimeError
		if errors.As(err, &trgt) {
			runtimeError(trgt)
		} else {
			panic(err)
		}
	}
	fmt.Println(stringify(value))
}

func checkNumberOperand(operator Token, operand any) error {
	_, ok := operand.(float64)
	if !ok {
		return RuntimeError{
			token: operator,
			msg:   "Operand must be a number",
		}
	}
	return nil
}

func checkNumberOperands(operator Token, left any, right any) error {
	_, lok := left.(float64)
	_, rok := right.(float64)

	if !lok || !rok {
		return RuntimeError{
			token: operator,
			msg:   "Operand must be a number",
		}
	}
	return nil
}

func evaluate(expr Expr) (any, error) {
	switch t := expr.(type) {
	case Binary:
		return evalBinary(t)
	case Grouping:
		return evalGrouping(t)
	case Literal:
		return evalLiteral(t), nil
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

func stringify(object any) string {
	return fmt.Sprintf("%v", object)
}

func evalBinary(expr Binary) (any, error) {
	left, err := evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case GREATER:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, fmt.Errorf("checking binary greater than: %w", err)
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, fmt.Errorf("checking binary greater than or equal: %w", err)
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, fmt.Errorf("checking binary less than: %w", err)
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, fmt.Errorf("checking binary less than or equal: %w", err)
		}
		return left.(float64) <= right.(float64), nil
	case MINUS:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, fmt.Errorf("checking binary subtraction: %w", err)
		}
		return left.(float64) - right.(float64), nil
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	case PLUS:
		// Pluss is a bit special because it works for
		// numbers and strings
		{
			l, lok := left.(float64)
			r, rok := right.(float64)
			if lok && rok {
				return l + r, nil
			}
		}
		{
			l, lok := left.(string)
			r, rok := right.(string)
			if lok && rok {
				return l + r, nil
			}
		}
		return nil, fmt.Errorf("checking plus (could be number or string): %w", RuntimeError{
			token: expr.operator,
			msg:   "Operands must be two numbers or two strings",
		})
	case SLASH:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, fmt.Errorf("checking binary division (SLASH): %w", err)
		}
		// Check if we are dividing by zero
		rval := right.(float64)
		if rval == 0.0 {
			return nil, RuntimeError{
				token: expr.operator,
				msg:   "divide by zero",
			}
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, fmt.Errorf("checking binary multiplication (STAR): %w", err)
		}
		return left.(float64) * right.(float64), nil
	}

	// Unreachable
	panic("eval binary: should never get here...")
	return nil, fmt.Errorf("eval binary: should never get here...")
}

func evalGrouping(expr Grouping) (any, error) {
	return evaluate(expr.expression)
}

func evalLiteral(expr Literal) any {
	return expr.value
}

func evalUnary(expr Unary) (any, error) {
	right, err := evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case BANG:
		return !isTruthy(right), nil
	case MINUS:
		if err := checkNumberOperand(expr.operator, right); err != nil {
			return nil, fmt.Errorf("checking minus operand: %w", err)
		}
		return -(right.(float64)), nil
	}

	// Unreachable
	panic("eval unary: should never get here...")
	return nil, fmt.Errorf("eval unary: should never get here...")
}
