package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	hadError        bool
	hadRuntimeError bool
)

// Making this global because it is a static field on the Lox class in the book
// https://craftinginterpreters.com/evaluating-expressions.html#running-the-interpreter
var interpreter = Interpreter{
	ENvironment: NewEnvironment(nil),
}

func lmain() {
	switch {
	case len(os.Args) > 2:
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	case len(os.Args) == 2:
		runFile(os.Args[1])
	default:
		runPrompt()
	}
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	run(string(data))

	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)

	fmt.Printf("> ")
	for s.Scan() {
		line := s.Text()
		run(line)
		hadError = false
		fmt.Printf("> ")
	}
}

func run(source string) {
	scanner := &Scanner{source: source}
	tokens := scanner.scanTokens()

	parser := NewParser(tokens)
	statements, err := parser.parse()
	if err != nil {
		panic(err)
	}

	// Stop if there was a syntax error.
	if hadError {
		return
	}

	// Uncomment to print the ast for debu
	// fmt.Println(ASTPrint(expression))

	interpreter.interpret(statements)
}

func loxlineerror(line int, message string) {
	loxreport(line, "", message)
}

func runtimeError(err RuntimeError) {
	fmt.Printf("RUNTIME ERROR: %v\n[line %v]\n", err, err.token.line)
	hadRuntimeError = true
}

func loxreport(line int, where, message string) {
	fmt.Printf("[line %d] Error%v: %v\n", line, where, message)
	hadError = true
}

func loxtokenerror(token Token, message string) {
	if token.tokenType == EOF {
		loxreport(token.line, " at end", message)
	} else {
		loxreport(token.line, " at '"+token.lexeme+"'", message)
	}
}
