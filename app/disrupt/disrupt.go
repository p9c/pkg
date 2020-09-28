package disrupt

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kardianos/osext"

	"github.com/p9c/pkg/app/slog"
)

var (
	Restart   bool // = true
	requested bool
	// Chan is used to receive SIGINT (Ctrl+C) signals.
	Chan chan os.Signal
	// Signals is the list of signals that cause the interrupt
	Signals = []os.Signal{os.Interrupt}
	// ShutdownRequestChan is a channel that can receive shutdown requests
	ShutdownRequestChan = make(chan struct{})
	// AddHandlerChan is used to add an interrupt handler to the list of
	// handlers to be invoked on SIGINT (Ctrl+C) signals.
	AddHandlerChan = make(chan func())
	// HandlersDone is closed after all interrupt handlers run the first time
	// an interrupt is signaled.
	HandlersDone = make(chan struct{})
)

func init() {
	// by importing this library the shutdown listener will be automatically started and by adding closures using
	// AddHandler you can be sure that short of sigkill or hardware failure that your application will run shutdown code
	// before it terminates
	go Listener()
}

// Receiver listens for interrupt signals, registers interrupt callbacks, and responds to custom shutdown signals as
// required
func Listener() {
	var interruptCallbacks []func()
	invokeCallbacks := func() {
		slog.Debug("running interrupt callbacks")
		// run handlers in LIFO order.
		for i := range interruptCallbacks {
			idx := len(interruptCallbacks) - 1 - i
			interruptCallbacks[idx]()
		}
		close(HandlersDone)
		slog.Debug("interrupt handlers finished")
		if Restart {
			slog.Debug("restarting")
			file, err := osext.Executable()
			if err != nil {
				slog.Error(err)
				return
			}
			err = syscall.Exec(file, os.Args, os.Environ())
			if err != nil {
				slog.Fatal(err)
			}
		}
		os.Exit(1)
	}
	for {
		select {
		case sig := <-Chan:
			// L.Printf("\r>>> received signal (%s)\n", sig)
			fmt.Println()
			slog.Debug("received interrupt signal", sig)
			requested = true
			invokeCallbacks()
			return
		case <-ShutdownRequestChan:
			slog.Warn("received shutdown request - shutting down...")
			requested = true
			invokeCallbacks()
			return
		case handler := <-AddHandlerChan:
			slog.Debug("adding handler")
			interruptCallbacks = append(interruptCallbacks, handler)
		}
	}
}

// AddHandler adds a handler to call when a SIGINT (Ctrl+C) is received.
func AddHandler(handler func()) {
	// Create the channel and start the main interrupt handler which invokes all other callbacks and exits if not
	// already done.
	if Chan == nil {
		Chan = make(chan os.Signal, 1)
		signal.Notify(Chan, Signals...)
		go Listener()
	}
	AddHandlerChan <- handler
}

// Request programatically requests a shutdown
func Request() {
	slog.Debug("interrupt requested")
	ShutdownRequestChan <- struct{}{}
}

// RequestRestart sets the reset flag and requests a restart
func RequestRestart() {
	Restart = true
	slog.Debug("requesting restart")
	Request()
}

// Requested returns true if an interrupt has been requested
func Requested() bool {
	return requested
}
