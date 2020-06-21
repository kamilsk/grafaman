package debugger

import (
	"net"
	"net/http"

	"github.com/pkg/errors"
	config "go.octolab.org/toolkit/config/http"
)

// Option defines function to configure debugger.
type Option func(*debugger) error

// WithBuiltinServer configures debugger by built-in HTTP server.
func WithBuiltinServer(config config.Server) Option {
	return func(debugger *debugger) error {
		listener, err := net.Listen("tcp", config.Address)
		if err != nil {
			return errors.Wrap(err, "debugger: listen tcp")
		}
		server := &http.Server{
			Addr:              listener.Addr().String(),
			Handler:           http.DefaultServeMux,
			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,
			MaxHeaderBytes:    config.MaxHeaderBytes,
		}
		return WithCustomListenerAndServer(listener, server)(debugger)
	}
}

// WithCustomListenerAndServer configures debugger by custom listener and server.
func WithCustomListenerAndServer(listener Listener, server Server) Option {
	return func(debugger *debugger) error {
		debugger.listener = listener
		debugger.server = server
		return nil
	}
}

// WithSpecificHost configures debugger by specific host.
func WithSpecificHost(host string) Option {
	return func(debugger *debugger) error {
		return WithBuiltinServer(config.Server{Address: host})(debugger)
	}
}
