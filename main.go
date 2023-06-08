package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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

func runFile(_ string) {
	fmt.Println("runFile() not implemented yet")
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)

	fmt.Printf("> ")
	for s.Scan() {
		line := s.Text()
		run(line)
		fmt.Printf("> ")
	}
}

func run(_ string) {
	fmt.Println("run() not implemented yet")
}
