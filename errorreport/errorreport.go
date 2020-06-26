package errorreport

import "fmt"

type ErrorReport struct {
	HadError bool
}

func (report *ErrorReport) Report(line int, where string, message string) {
	report.HadError = true
	fmt.Printf("[line %d] Error %v: %v\n", line, where, message)
}
