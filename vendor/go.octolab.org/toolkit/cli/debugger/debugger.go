package debugger

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"go.octolab.org/safe"
)

// Must returns new configured debugger or panic if an error occurred.
func Must(options ...Option) *debugger {
	debugger, err := New(options...)
	if err != nil {
		panic(err)
	}
	return debugger
}

// New returns new configured debugger or an error if something went wrong.
func New(options ...Option) (*debugger, error) {
	debugger := new(debugger)
	for i, configure := range options {
		if err := configure(debugger); err != nil {
			return nil, errors.Wrapf(err, "debugger: configure by %d option", i+1)
		}
	}
	if debugger.listener == nil || debugger.server == nil {
		return nil, errors.New("debugger: without listener or server")
	}
	return debugger, nil
}

type debugger struct {
	debug    sync.Once
	listener Listener
	server   Server
}

// Debug runs debugger only once and returns the fact of success run.
func (debugger *debugger) Debug(logger func(error), shutdown ...func()) (string, bool) {
	var status uint32
	debugger.debug.Do(func() {
		for _, fn := range shutdown {
			debugger.server.RegisterOnShutdown(fn)
		}
		go safe.Do(func() error {
			if err := debugger.server.Serve(debugger.listener); err != http.ErrServerClosed {
				return errors.Wrap(err, "debugger: serve listener")
			}
			return nil
		}, logger)
		atomic.CompareAndSwapUint32(&status, 0, 1)
	})
	return debugger.listener.Addr().String(), atomic.CompareAndSwapUint32(&status, 1, 0)
}

// Stop tries to stop debugger if it runs.
func (debugger *debugger) Stop(ctx context.Context) error {
	return errors.Wrap(debugger.server.Shutdown(ctx), "debugger: shutdown server")
}
