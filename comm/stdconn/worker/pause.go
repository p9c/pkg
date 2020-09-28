// +build !windows

package worker

import (
	"syscall"

	"github.com/p9c/pkg/app/slog"
)

func (w *Worker) Pause() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGSTOP); !slog.Check(err) {
		slog.Debug("paused")
	}
	return
}
func (w *Worker) Resume() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGCONT); !slog.Check(err) {
		slog.Debug("resumed")
	}
	return
}
