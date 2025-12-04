---
id: 12
database_id: 625478326
node_id: MDU6SXNzdWU2MjU0NzgzMjY=
status: closed
title: "critical bug with nil pointer"
labels: ["bug","help wanted"]
url: https://github.com/kamilsk/grafaman/issues/12
created_at: 2020-05-27T08:12:29Z
updated_at: 2020-05-28T06:04:34Z
---

# critical bug with nil pointer

```bash
$ grafaman version
recovered: runtime error: invalid memory address or nil pointer dereference
---
unexpected panic occurred
go.octolab.org/safe.Do.func2
	/home/travis/gopath/pkg/mod/go.octolab.org@v0.0.26/safe/do.go:27
runtime.gopanic
	/home/travis/.gimme/versions/go1.14.3.linux.amd64/src/runtime/panic.go:969
runtime.panicmem
	/home/travis/.gimme/versions/go1.14.3.linux.amd64/src/runtime/panic.go:212
runtime.sigpanic
	/home/travis/.gimme/versions/go1.14.3.linux.amd64/src/runtime/signal_unix.go:695
github.com/spf13/viper.(*Viper).AllKeys
	/home/travis/gopath/pkg/mod/github.com/spf13/viper@v1.6.3/viper.go:1812
github.com/spf13/viper.(*Viper).AllSettings
	/home/travis/gopath/pkg/mod/github.com/spf13/viper@v1.6.3/viper.go:1892
github.com/kamilsk/grafaman/internal/cmd.New.func1
	/home/travis/gopath/src/github.com/kamilsk/grafaman/internal/cmd/root.go:40
github.com/spf13/cobra.(*Command).execute
	/home/travis/gopath/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:821
github.com/spf13/cobra.(*Command).ExecuteC
	/home/travis/gopath/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:950
github.com/spf13/cobra.(*Command).Execute
	/home/travis/gopath/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:887
github.com/spf13/cobra.(*Command).ExecuteContext
	/home/travis/gopath/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:880
main.main.func1
	/home/travis/gopath/src/github.com/kamilsk/grafaman/main.go:46
go.octolab.org/safe.Do
	/home/travis/gopath/pkg/mod/go.octolab.org@v0.0.26/safe/do.go:30
main.main
	/home/travis/gopath/src/github.com/kamilsk/grafaman/main.go:46
runtime.main
	/home/travis/.gimme/versions/go1.14.3.linux.amd64/src/runtime/proc.go:203
runtime.goexit
	/home/travis/.gimme/versions/go1.14.3.linux.amd64/src/runtime/asm_amd64.s:1373
```
