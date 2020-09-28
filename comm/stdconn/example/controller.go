package main

import (
	"github.com/p9c/pkg/app/slog"
	"github.com/p9c/pkg/comm/stdconn/example/hello/hello"
	"github.com/p9c/pkg/comm/stdconn/worker"
)

func main() {
	slog.Info("starting up example controller")
	cmd := worker.Spawn("go", "run", "hello/worker.go")
	client := hello.NewClient(cmd.StdConn)
	slog.Info("calling Hello.Say with 'worker'")
	slog.Info("reply:", client.Say("worker"))
	slog.Info("calling Hello.Bye")
	slog.Info("reply:", client.Bye())
}
