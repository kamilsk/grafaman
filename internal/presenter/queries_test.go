package presenter_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kamilsk/grafaman/internal/model"
	. "github.com/kamilsk/grafaman/internal/presenter"
)

func TestPrinter_PrintQueries(t *testing.T) {
	queries := model.Queries{
		"metric.a.*",
		"metric.b.*",
		"metric.c.*",
	}

	tests := map[string]struct {
		output interface {
			io.Writer
			fmt.Stringer
		}
		format string
		prefix string
		assert func(require.TestingT, error, string)
	}{
		"default": {
			output: bytes.NewBuffer(nil),
			format: (*Printer)(nil).DefaultFormat(),
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.default.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
		"compact": {
			output: bytes.NewBuffer(nil),
			format: "compact",
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.compact.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
		"compact-lite": {
			output: bytes.NewBuffer(nil),
			format: "compact-lite",
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.compact-lite.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
		"markdown": {
			output: bytes.NewBuffer(nil),
			format: "markdown",
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.markdown.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
		"rounded": {
			output: bytes.NewBuffer(nil),
			format: "rounded",
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.rounded.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
		"unicode": {
			output: bytes.NewBuffer(nil),
			format: "unicode",
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.unicode.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
		"json": {
			output: bytes.NewBuffer(nil),
			format: "json",
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.json.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
		"tsv": {
			output: bytes.NewBuffer(nil),
			format: "tsv",
			assert: func(t require.TestingT, err error, output string) {
				require.NoError(t, err)

				file := "testdata/queries.tsv.txt"
				if *update {
					require.NoError(t, ioutil.WriteFile(file, []byte(output), 0644))
				}

				golden, err := ioutil.ReadFile(file)
				assert.NoError(t, err)
				assert.Equal(t, string(golden), output)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			printer := new(Printer).SetOutput(test.output)
			printer.SetPrefix(test.prefix)
			require.NoError(t, printer.SetFormat(test.format))

			test.assert(t, printer.PrintQueries(queries), test.output.String())
		})
	}

	t.Run("fs unhealthy", func(t *testing.T) {
		printer := new(Printer).SetOutput(new(unhealthy))
		require.NoError(t, printer.SetFormat("tsv"))

		assert.Error(t, printer.PrintQueries(queries))
	})
}
