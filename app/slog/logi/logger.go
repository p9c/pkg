// Package logi is a logger interface designed to allow easy extension with log printer functions and complete
// reimplementation as required
package logi

type Level int
type Printer func(lvl Level, txt interface{})

// Logger is an interface that defines the set of operations relevant to a logger. This interface should be a stdlib
// interface really because logging is so useful for debugging.
type Logger interface {
	// The following are printers that print at or below the given level from the constants above
	Fatal(txt ...interface{})
	Error(txt ...interface{})
	Warn(txt ...interface{})
	Info(txt ...interface{})
	Debug(txt ...interface{})
	Trace(txt ...interface{})
	Fatalf(format string, txt ...interface{})
	Errorf(format string, txt ...interface{})
	Warnf(format string, txt ...interface{})
	Infof(format string, txt ...interface{})
	Debugf(format string, txt ...interface{})
	Tracef(format string, txt ...interface{})
	Errors(txt interface{})
	Warns(txt interface{})
	Infos(txt interface{})
	Debugs(txt interface{})
	Traces(txt interface{})

	// Check prints at error level if the error was not nil and returns true
	Check(err error) bool
	// SetPrinter enables loading a printer function to enable networked, piped, etc outputs
	SetPrinter(fn Printer)
	// SetLevel sets the error level, anything lower will not call the printer function
	SetLevel(level Level)
}

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	BOUNDARY
)

var LevelCodes = []string{
	"TRC", "DBG", "INF", "WRN", "ERR", "FTL",
}
