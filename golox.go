package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/parser"
	"github.com/maleksiuk/golox/scanner"
	"github.com/maleksiuk/golox/tools"
)

func main() {
	args := os.Args[1:]

	argCount := len(args)
	switch {
	case argCount > 1:
		fmt.Println("Usage: golox [script]")
	case argCount == 1:
		err := runFile(args[0])
		if err != nil {
			os.Exit(1)
		}
	default:
		runPrompt()
	}
}

func runFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer file.Close()

	errorReport := errorreport.ErrorReport{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		run(scanner.Text(), &errorReport)

		if errorReport.HadError {
			return errors.New("Scanner error")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func runPrompt() {
	errorReport := errorreport.ErrorReport{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		run(scanner.Text(), &errorReport)
		errorReport.HadError = false
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func run(line string, errorReport *errorreport.ErrorReport) {
	tokens := scanner.ScanTokens(line, errorReport)
	expr := parser.Parse(tokens, errorReport)
	str := tools.PrintAst(expr)
	fmt.Println(str)
	for _, token := range tokens {
		fmt.Printf("token: %v\n", token)
	}
}
