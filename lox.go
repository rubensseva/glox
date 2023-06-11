package main

import (
	"bufio"
	"fmt"
	"os"
)

var hadError bool

func lmain() {
	switch {
	case len(os.Args) > 2:
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	case len(os.Args) == 2:
		runFile(os.Args[0])
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

	// For now, just print the tokens
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func loxerror(line int, message string) {
	loxreport(line, "", message)
}
func loxreport(line int, where, message string) {
	fmt.Printf("[line %d] Error%v: %v\n", line, where, message)
	hadError = true
}
