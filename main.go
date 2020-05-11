package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.octolab.org/errors"
	"go.octolab.org/safe"
	"go.octolab.org/toolkit/cli/cobra"
	"go.octolab.org/unsafe"

	"github.com/kamilsk/grafaman/internal/cmd"
)

const unknown = "unknown"

var (
	commit  = unknown
	date    = unknown
	version = "dev"
	exit    = os.Exit
	stderr  = io.Writer(os.Stderr)
	stdout  = io.Writer(os.Stdout)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	root := cmd.New()
	root.SetErr(stderr)
	root.SetOut(stdout)
	root.AddCommand(
		cobra.NewCompletionCommand(),
		cobra.NewVersionCommand(version, date, commit),
	)
	safe.Do(func() error { return root.ExecuteContext(ctx) }, shutdown)
}

func shutdown(err error) {
	if recovered, is := errors.Unwrap(err).(errors.Recovered); is {
		unsafe.DoSilent(fmt.Fprintf(stderr, "recovered: %+v\n", recovered.Cause()))
		unsafe.DoSilent(fmt.Fprintln(stderr, "---"))
		unsafe.DoSilent(fmt.Fprintf(stderr, "%+v\n", err))
	}
	exit(1)
}
