package simple

import (
	"errors"
	"testing"

	"github.com/p9c/pkg/app/slog/logi"
)

func TestSimpleLogger(t *testing.T) {
	l := NewSimpleLogger()
	l.Panic = false
	var lvl logi.Level
	for ; lvl < logi.BOUNDARY; lvl++ {
		l.SetLevel(lvl)
		t.Log("testing level", logi.LevelCodes[lvl])
		l.Trace("trace level")
		l.Debug("debug level")
		l.Info("info level")
		l.Warn("warn level")
		l.Error("error level")
		l.Error(l.Check(nil))
		l.Error(l.Check(errors.New("check test")))
		l.Fatal("fatal level")
	}
}
