package graphite

import "net/http"

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test

// A Client defines the basic HTTP client interface.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type ProgressListener interface {
	OnStepDone()
	OnStepQueued()
}
