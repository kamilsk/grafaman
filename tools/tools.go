// Code generated by github.com/kamilsk/egg. DO NOT EDIT.

// +build tools

package tools

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/kyoh86/looppointer/cmd/looppointer"
	_ "golang.org/x/tools/cmd/goimports"
)

//go:generate go install github.com/golang/mock/mockgen
//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install github.com/kyoh86/looppointer/cmd/looppointer
//go:generate go install golang.org/x/tools/cmd/goimports
