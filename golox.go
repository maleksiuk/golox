package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/maleksiuk/golox/errorreport"
	"github.com/maleksiuk/golox/interpreter"
	"github.com/maleksiuk/golox/parser"
	"github.com/maleksiuk/golox/scanner"
)

func main() {
	args := os.Args[1:]
	argCount := len(args)

	i := interpreter.NewInterpreter()

	switch {
	case argCount > 1:
		fmt.Println("Usage: golox [script]")
	case argCount == 1:
		err := runFile(i, args[0])
		if err != nil {
			os.Exit(1)
		}
	default:
		runPrompt(i)
	}
}

func runFile(i interpreter.Interpreter, path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer file.Close()

	errorReport := errorreport.ErrorReport{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		run(i, scanner.Text(), &errorReport)

		if errorReport.HadError {
			return errors.New("Scanner error")
		}
		if errorReport.HadRuntimeError {
			return errors.New("Runtime error")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func runPrompt(i interpreter.Interpreter) {
	errorReport := errorreport.ErrorReport{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		run(i, scanner.Text(), &errorReport)
		errorReport.HadError = false
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func run(i interpreter.Interpreter, line string, errorReport *errorreport.ErrorReport) {
	tokens := scanner.ScanTokens(line, errorReport)
	statements := parser.Parse(tokens, errorReport)

	i.Interpret(statements, errorReport)
}
