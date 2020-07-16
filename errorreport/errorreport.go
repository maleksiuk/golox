package errorreport

import "fmt"

type ErrorReport struct {
	HadError        bool
	HadRuntimeError bool
}

func (report *ErrorReport) Report(line int, where string, message string) {
	report.HadError = true
	fmt.Printf("[line %d] Error %v: %v\n", line, where, message)
}

func (report *ErrorReport) ReportRuntimeError(line int, message string) {
	report.HadRuntimeError = true
	fmt.Printf("[line %d] Runtime error: %v\n", line, message)
}
