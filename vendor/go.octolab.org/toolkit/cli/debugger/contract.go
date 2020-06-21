package debugger

import (
	"context"
	"net"
)

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test

// A Server is a generic server to listen a network protocol.
// It is compatible with net/http.Server.
type Server interface {
	// Serve accepts incoming connections on the Listener and serves them.
	Serve(net.Listener) error
	// RegisterOnShutdown registers a function to call on Shutdown.
	RegisterOnShutdown(func())
	// Shutdown tries to do a graceful shutdown.
	Shutdown(context.Context) error
}

// A Listener is a generic network listener for stream-oriented protocols.
type Listener interface {
	net.Listener
}
