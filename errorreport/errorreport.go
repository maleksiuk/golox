package errorreport

type ErrorReport struct {
	HadError        bool
	HadRuntimeError bool
	Printer         Printer
}

func NewErrorReport() ErrorReport {
	return ErrorReport{Printer: consolePrinter{}}
}

func (report *ErrorReport) Report(line int, where string, message string) {
	report.HadError = true
	report.Printer.Printf("[line %d] Error %v: %v\n", line, where, message)
}

func (report *ErrorReport) ReportRuntimeError(line int, message string) {
	report.HadRuntimeError = true
	report.Printer.Printf("[line %d] Runtime error: %v\n", line, message)
}
