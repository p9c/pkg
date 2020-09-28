package worker

import "github.com/p9c/pkg/app/slog"

func (w *Worker) Pause() (err error) {
	slog.Debug("windows can't pause processes")
	return
}
func (w *Worker) Resume() (err error) {
	slog.Debug("windows can't pause so can't resume processes")
	return
}
