package main

import (
	"github.com/p9c/pkg/app/disrupt"
	"github.com/p9c/pkg/app/slog"
)

func main() {
	disrupt.AddHandler(func() {
		slog.Warn("IT'S THE END OF THE WORLD!")
	})
	<-disrupt.HandlersDone
}
