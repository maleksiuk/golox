package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	argCount := len(args)
	switch {
	case argCount > 1:
		fmt.Println("Usage: golox [script]")
	case argCount == 1:
		runFile(args[0])
	default:
		runPrompt()
	}
}

func runFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		run(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		run(scanner.Text())
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func run(line string) {
	fmt.Println(line)
}
