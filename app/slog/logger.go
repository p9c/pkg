// Package log creates a neat short access to logger functions for a simple app-wide namespace (package global, this
// package need only be imported) for use of slog/simple logger implementation
package slog

import (
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"

	"github.com/p9c/pkg/app/slog/logi"
	"github.com/p9c/pkg/app/slog/simple"
)

// a simple, convenient logger with configuration of log level through an OS environment variable
func init() {
	l := simple.NewSimpleLogger()
	ls := os.Getenv("LOGLEVEL")
	ll, err := strconv.ParseInt(ls, 10, 8)
	logLevel := logi.Level(ll)
	if ll > int64(logi.BOUNDARY) || err != nil {
		logLevel = logi.INFO
	}
	SetLevel = l.SetLevel
	// Debug("set log level", logi.LevelCodes[logLevel])
	SetPrinter = l.SetPrinter
	Trace = l.Trace
	Debug = l.Debug
	Info = l.Info
	Warn = l.Warn
	Error = l.Error
	Check = l.Check
	Fatal = l.Fatal
	Tracef = l.Tracef
	Debugf = l.Debugf
	Infof = l.Infof
	Warnf = l.Warnf
	Errorf = l.Errorf
	Fatalf = l.Fatalf
	SetLevel(logLevel)
	Traces = func(txt interface{}) {
		l.Trace(spew.Sdump(txt))
	}
	Debugs = func(txt interface{}) {
		l.Debug(spew.Sdump(txt))
	}
	Infos = func(txt interface{}) {
		l.Info(spew.Sdump(txt))
	}
	Warns = func(txt interface{}) {
		l.Warn(spew.Sdump(txt))
	}
	Errors = func(txt interface{}) {
		l.Error(spew.Sdump(txt))
	}
	Fatals = func(txt interface{}) {
		l.Fatal(spew.Sdump(txt))
	}
}

// The following are slots that can be loaded for an app-wide logger
var (
	Fatal      func(txt ...interface{})
	Error      func(txt ...interface{})
	Warn       func(txt ...interface{})
	Info       func(txt ...interface{})
	Debug      func(txt ...interface{})
	Trace      func(txt ...interface{})
	Fatalf     func(format string, txt ...interface{})
	Errorf     func(format string, txt ...interface{})
	Warnf      func(format string, txt ...interface{})
	Infof      func(format string, txt ...interface{})
	Debugf     func(format string, txt ...interface{})
	Tracef     func(format string, txt ...interface{})
	Fatals     func(txt interface{})
	Errors     func(txt interface{})
	Warns      func(txt interface{})
	Infos      func(txt interface{})
	Debugs     func(txt interface{})
	Traces     func(txt interface{})
	Check      func(err error) bool
	SetPrinter func(fn logi.Printer)
	SetLevel   func(level logi.Level)
)
