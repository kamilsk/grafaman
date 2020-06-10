module github.com/kamilsk/grafaman/tools

go 1.13

require (
	github.com/golang/mock v1.4.3
	github.com/golangci/golangci-lint v1.27.0
	github.com/kamilsk/egg v0.0.16
	golang.org/x/exp v0.0.0-20200513190911-00229845015e
	golang.org/x/tools v0.3.3
)

replace github.com/izumin5210/gex => github.com/kamilsk/gex v0.6.0-e4

replace golang.org/x/tools => github.com/kamilsk/go-tools v0.0.3
