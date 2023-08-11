package main

import (
	"errors"
	"fmt"
)

type ReturnHack struct {
	value any
}

type Interpreter struct {
	globals *Environment
	// Using ENvironment name on purpose to closely match the book
	ENvironment *Environment
}

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{
		globals: NewEnvironment(nil),
	}
	interpreter.ENvironment = interpreter.globals

	interpreter.globals.define("clock", &Clock{})

	return interpreter
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			var trgt RuntimeError
			if errors.As(err, &trgt) {
				runtimeError(trgt)
				break
			} else {
				panic(err)
			}
		}
	}
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

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	switch t := expr.(type) {
	case Binary:
		return i.visitBinaryExpr(t)
	case Grouping:
		return i.visitGroupingExpr(t)
	case Literal:
		return i.visitLiteralExpr(t), nil
	case Unary:
		return i.visitUnaryExpr(t)
	case Variable:
		return i.visitVariableExpr(t)
	case Logical:
		return i.visitLogicalExpr(t)
	case Assign:
		return i.visitAssignExpr(t)
	case Call:
		return i.visitCallExpr(t)
	default:
		panic(fmt.Sprintf("eval: unknown type %T: %v", expr, t))
	}
}

func (i *Interpreter) execute(stmt Stmt) error {
	switch t := stmt.(type) {
	case PrintStmt:
		return i.visitPrintStmt(t)
	case VarStmt:
		return i.visitVarStmt(t)
	case ExpressionStmt:
		i.visitExpressionStmt(t)
		return nil
	case BlockStmt:
		i.visitBlockStmt(t)
		return nil
	case IfStmt:
		if err := i.visitIfStmt(t); err != nil {
			return fmt.Errorf("visiting if statemenet: %w", err)
		}
		return nil
	case WhileStmt:
		if err := i.visitWhileStmt(t); err != nil {
			return fmt.Errorf("visiting wihle statement: %w", err)
		}
		return nil
	case FunctionStmt:
		i.visitFunctionStmt(t)
		return nil
	case ReturnStmt:
		if err := i.visitReturnStmt(t); err != nil {
			return err
		}
		return nil
	default:
		panic(fmt.Sprintf("executing: unknown type %T: %v", stmt, t))
	}
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	// https://craftinginterpreters.com/statements-and-state.html#block-syntax-and-semantics
	previous := i.ENvironment
	defer func() {
		i.ENvironment = previous
	}()

	i.ENvironment = environment
	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) visitBlockStmt(stmt BlockStmt) {
	if err := i.executeBlock(stmt.statements, NewEnvironment(i.ENvironment)); err != nil {
		panic(err)
	}
}

func (i *Interpreter) visitExpressionStmt(stmt ExpressionStmt) {
	_, err := i.evaluate(stmt.expression)
	if err != nil {
		panic(err)
	}
}

func (i *Interpreter) visitFunctionStmt(stmt FunctionStmt) {
	function := NewLoxFunction(stmt, i.ENvironment)
	i.ENvironment.define(stmt.name.lexeme, function)
}

// https://craftinginterpreters.com/control-flow.html#conditional-execution
func (i *Interpreter) visitIfStmt(stmt IfStmt) error {
	evres, err := i.evaluate(stmt.condition)
	if err != nil {
		return fmt.Errorf("evaluating condition: %w", err)
	}
	if isTruthy(evres) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
	return nil
}

func (i Interpreter) visitPrintStmt(stmt PrintStmt) error {
	var value any
	var err error
	value, err = i.evaluate(stmt.expression)
	if err != nil {
		return fmt.Errorf("evaluating print printstmt expression: %w", err)
	}
	fmt.Println(stringify(value))
	return nil
}

func (i Interpreter) visitReturnStmt(stmt ReturnStmt) error {
	var value any = nil
	if stmt.value != nil {
		tmp, err := i.evaluate(stmt.value)
		if err != nil {
			return err
		}
		value = tmp
	}

	// HACK: Absolutely massive hack inoming!
	// We are using panic() as control flow here
	// The book uses Java exceptions, so this is an attempt
	// to emulate that as close as possible
	panic(ReturnHack{value: value})
}

func (i *Interpreter) visitVarStmt(stmt VarStmt) error {
	var value any = nil
	if stmt.initializer != nil {
		v, err := i.evaluate(stmt.initializer)
		if err != nil {
			return fmt.Errorf("evaluating initializer: %w", err)
		}
		value = v
	}
	i.ENvironment.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) visitWhileStmt(stmt WhileStmt) error {
	for {
		condEvald, err := i.evaluate(stmt.condition)
		if err != nil {
			return fmt.Errorf("evaluating stmt condition in while loop: %w", err)
		}
		if !isTruthy(condEvald) {
			break
		}
		i.execute(stmt.body)
	}
	return nil
}

func (i *Interpreter) visitAssignExpr(expr Assign) (any, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, fmt.Errorf("evaluating assignment expression: %w", err)
	}
	i.ENvironment.assign(expr.name, value)
	return value, nil
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

func (i *Interpreter) visitBinaryExpr(expr Binary) (any, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
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
			msg: fmt.Sprintf("Operands must be two numbers or two strings, got %[1]v %[1]T, %[2]v %[2]T",
				left, right),
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

func (i *Interpreter) visitCallExpr(expr Call) (any, error) {
	callee, err := i.evaluate(expr.callee)
	if err != nil {
		return nil, fmt.Errorf("evaluate(): %w", err)
	}

	var arguments []any
	for _, argument := range expr.arguments {
		e, err := i.evaluate(argument)
		if err != nil {
			return nil, fmt.Errorf("evaluate(): %w", err)
		}
		arguments = append(
			arguments,
			e,
		)
	}

	// TODO: Mabe type cast instead?
	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, fmt.Errorf("callee was not a LoxCallable: %[1]v, %[1]T", callee)
	}
	if len(arguments) != function.Arity() {
		return nil, RuntimeError{
			token: expr.paren,
			msg: fmt.Sprintf("Expected %d arguments but got %d.",
				function.Arity(), len(arguments)),
		}
	}

	res, err := function.Call(i, arguments)
	if err != nil {
		return nil, fmt.Errorf("LoxCallable.call(): %w", err)
	}
	return res, nil
}

func (i *Interpreter) visitGroupingExpr(expr Grouping) (any, error) {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitLiteralExpr(expr Literal) any {
	return expr.value
}

func (i *Interpreter) visitLogicalExpr(expr Logical) (any, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, fmt.Errorf("evaluating left expr of logical expr: %w", err)
	}

	if expr.operator.tokenType == OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		// I think this branch means that we assume expr.operator.tokenType == AND ?
		// See chapter 9.3
		if !isTruthy(left) {
			return left, nil
		}
	}

	res, err := i.evaluate(expr.right)
	if err != nil {
		return nil, fmt.Errorf("evaluating right expr: %w", err)
	}
	return res, nil
}

func (i *Interpreter) visitUnaryExpr(expr Unary) (any, error) {
	right, err := i.evaluate(expr.right)
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

func (i *Interpreter) visitVariableExpr(expr Variable) (any, error) {
	return i.ENvironment.get(expr.name)
}
