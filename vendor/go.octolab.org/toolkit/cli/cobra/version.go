package cobra

import (
	"fmt"
	"runtime"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// NewVersionCommand returns a command that helps to build version info.
//
//  $ cli version
//  cli:
//    version     : 1.0.0
//    build date  : 2019-07-17T12:44:00Z
//    git hash    : 4f8c7f4
//    go version  : go1.12.7
//    go compiler : gc
//    platform    : darwin/amd64
//    features    : featureA=true, featureB=false
//
func NewVersionCommand(release, date, hash string, features ...Feature) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show application version",
		Long:  "Show application version.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return version.Execute(cmd.OutOrStdout(), struct {
				Name       string
				Version    string
				BuildDate  string
				GitHash    string
				GoVersion  string
				GoCompiler string
				Platform   string
				Features   fmt.Stringer
			}{
				Name:       root(cmd).Name(),
				Version:    release,
				BuildDate:  date,
				GitHash:    hash,
				GoVersion:  runtime.Version(),
				GoCompiler: runtime.Compiler,
				Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
				Features:   Features(features),
			})
		},
	}
}

// Feature describe a feature.
type Feature struct {
	Name    string
	Enabled bool
}

// String returns a string representation of the feature.
func (feature Feature) String() string {
	return fmt.Sprintf("%s=%v", feature.Name, feature.Enabled)
}

// Features defines a list of features.
type Features []Feature

// String returns a string representation of the feature list.
func (features Features) String() string {
	if len(features) == 0 {
		return "-"
	}
	list := make([]string, 0, len(features))
	for _, feature := range features {
		list = append(list, feature.String())
	}
	return strings.Join(list, ", ")
}

var version = template.Must(template.New("version").Parse(`{{.Name}}:
  version     : {{.Version}}
  build date  : {{.BuildDate}}
  git hash    : {{.GitHash}}
  go version  : {{.GoVersion}}
  go compiler : {{.GoCompiler}}
  platform    : {{.Platform}}
  features    : {{.Features}}
`))
