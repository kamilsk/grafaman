package http

import "time"

// Server contains configuration for a common HTTP server.
type Server struct {
	Address           string        `json:"address"             yaml:"address"`
	ReadTimeout       time.Duration `json:"read-timeout"        yaml:"read-timeout"`
	ReadHeaderTimeout time.Duration `json:"read-header-timeout" yaml:"read-header-timeout"`
	WriteTimeout      time.Duration `json:"write-timeout"       yaml:"write-timeout"`
	IdleTimeout       time.Duration `json:"idle-timeout"        yaml:"idle-timeout"`
	MaxHeaderBytes    int           `json:"max-header-bytes"    yaml:"max-header-bytes"`
}
