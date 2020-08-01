package presenter

import (
	"io"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"
)

// A Printer provides functionality to output data.
type Printer struct {
	format string
	prefix string
	output io.Writer
}

// SetFormat sets output format, e.g. json or tsv.
func (printer *Printer) SetFormat(format string) error {
	var present bool
	for _, supported := range formats {
		if format == supported {
			present = true
			break
		}
	}
	if !present {
		return errors.Errorf("presenter: invalid format %q, only %v are supported", format, formats)
	}
	printer.format = format
	return nil
}

// SetPrefix sets prefix to trim it from metric names.
func (printer *Printer) SetPrefix(prefix string) {
	printer.prefix = prefix
}

// SetOutput sets output.
func (printer *Printer) SetOutput(output io.Writer) *Printer {
	printer.output = output
	return printer
}

const (
	DefaultFormat     = "default"
	formatCompact     = "compact"
	formatCompactLite = "compact-lite"
	formatMarkdown    = "markdown"
	formatRounded     = "rounded"
	formatUnicode     = "unicode"
	formatJSON        = "json"
	formatTSV         = "tsv"
)

var (
	formats = []string{
		DefaultFormat,
		formatCompact,
		formatCompactLite,
		formatMarkdown,
		formatRounded,
		formatUnicode,
		formatJSON,
		formatTSV,
	}
	styles = map[string]*simpletable.Style{
		DefaultFormat:     simpletable.StyleDefault,
		formatCompact:     simpletable.StyleCompact,
		formatCompactLite: simpletable.StyleCompactLite,
		formatMarkdown:    simpletable.StyleMarkdown,
		formatRounded:     simpletable.StyleRounded,
		formatUnicode:     simpletable.StyleUnicode,
	}
)
