// Package simple is an implementation of the logi.Logger interface that prints log entries to stdout with level code,
// compact since startup time.Duration and appends the code location of the call to the logger at the end of the
// log text.
package simple

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/p9c/pkg/app/slog/logi"
)

type Simple struct {
	startupTime time.Time
	Printer     logi.Printer
	Level       logi.Level
	Panic       bool
}

var _ logi.Logger = (*Simple)(nil)

func (s *Simple) GetSimplePrinter() logi.Printer {
	return func(lvl logi.Level, txt interface{}) {
		if s.Level <= lvl {
			_, loc, line, _ := runtime.Caller(2)
			_, _ = fmt.Fprintln(os.Stderr,
				logi.LevelCodes[lvl],
				time.Now().Sub(s.startupTime).Round(time.Microsecond),
				txt, loc+":"+fmt.Sprint(line))
		}
	}
}

func NewSimpleLogger() (s *Simple) {
	s = &Simple{
		startupTime: time.Now(),
		Level:       logi.TRACE,
	}
	s.Printer = s.GetSimplePrinter()
	return
}

func (s *Simple) Fatal(txt ...interface{}) {
	s.Printer(logi.FATAL, txt)
	if s.Panic {
		panic("fatal error, printing stack trace and terminating")
	}
}

func (s *Simple) Error(txt ...interface{}) {
	s.Printer(logi.ERROR, txt)
}

func (s *Simple) Warn(txt ...interface{}) {
	s.Printer(logi.WARN, txt)
}

func (s *Simple) Info(txt ...interface{}) {
	s.Printer(logi.INFO, txt)
}

func (s *Simple) Debug(txt ...interface{}) {
	s.Printer(logi.DEBUG, txt)
}

func (s *Simple) Trace(txt ...interface{}) {
	s.Printer(logi.TRACE, txt)
}

func (s *Simple) Fatalf(format string, txt ...interface{}) {
	s.Printer(logi.FATAL, fmt.Sprintf(format, txt...))
	if s.Panic {
		panic("fatal error, printing stack trace and terminating")
	}
}

func (s *Simple) Errorf(format string, txt ...interface{}) {
	s.Printer(logi.ERROR, fmt.Sprintf(format, txt...))
}

func (s *Simple) Warnf(format string, txt ...interface{}) {
	s.Printer(logi.WARN, fmt.Sprintf(format, txt...))
}

func (s *Simple) Infof(format string, txt ...interface{}) {
	s.Printer(logi.INFO, fmt.Sprintf(format, txt...))
}

func (s *Simple) Debugf(format string, txt ...interface{}) {
	s.Printer(logi.DEBUG, fmt.Sprintf(format, txt...))
}

func (s *Simple) Tracef(format string, txt ...interface{}) {
	s.Printer(logi.TRACE, fmt.Sprintf(format, txt...))
}

func (s *Simple) Errors(txt interface{}) {
	s.Error(spew.Sdump(txt))
}

func (s *Simple) Warns(txt interface{}) {
	s.Warn(spew.Sdump(txt))
}

func (s *Simple) Infos(txt interface{}) {
	s.Info(spew.Sdump(txt))
}

func (s *Simple) Debugs(txt interface{}) {
	s.Debug(spew.Sdump(txt))
}

func (s *Simple) Traces(txt interface{}) {
	s.Trace(spew.Sdump(txt))
}

func (s *Simple) Check(err error) (errs bool) {
	if err != nil {
		s.Printer(logi.DEBUG, err)
		errs = true
	}
	return
}

func (s *Simple) SetPrinter(f logi.Printer) {
	s.Printer = f
}

func (s *Simple) SetLevel(level logi.Level) {
	if level < logi.BOUNDARY {
		s.Level = level
	}
}
