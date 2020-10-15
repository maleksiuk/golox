package errorreport

import "fmt"

type Printer interface {
	Printf(format string, a ...interface{}) (n int, err error)
}

type consolePrinter struct {
}

func (printer consolePrinter) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

// TODO: can this live in a test helper?
type MockPrinter struct {
	strings []string
}

func NewMockPrinter() *MockPrinter {
	return &MockPrinter{strings: make([]string, 0, 50)}
}

func (printer *MockPrinter) Printf(format string, a ...interface{}) (n int, err error) {
	str := fmt.Sprintf(format, a...)
	printer.strings = append(printer.strings, str)
	return 0, nil
}

func (printer *MockPrinter) GetStrings() []string {
	return printer.strings
}
