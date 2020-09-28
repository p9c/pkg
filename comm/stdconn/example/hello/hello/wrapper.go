package hello

import (
	"io"
	"net/rpc"

	"github.com/p9c/pkg/app/slog"
)

type Client struct {
	*rpc.Client
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}

}

func (h *Client) Say(name string) (reply string) {
	if err := h.Call("Hello.Say", "worker", &reply); slog.Check(err) {
		return "error: " + err.Error()
	}
	return
}

func (h *Client) Bye() (reply string) {
	if err := h.Call("Hello.Bye", 1, &reply); slog.Check(err) {
		return "error: " + err.Error()
	}
	return
}
