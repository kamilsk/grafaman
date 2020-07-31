package grafana

import "net/http"

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test

// Client defines HTTP client interface.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
