package presenter

import (
	"io"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"
)

type Printer struct {
	format string
	output io.Writer
}

func (printer *Printer) DefaultFormat() string {
	return formatDefault
}

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

func (printer *Printer) SetOutput(output io.Writer) *Printer {
	printer.output = output
	return printer
}

const (
	formatDefault     = "default"
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
		formatDefault,
		formatCompact,
		formatCompactLite,
		formatMarkdown,
		formatRounded,
		formatUnicode,
		formatJSON,
		formatTSV,
	}
	styles = map[string]*simpletable.Style{
		formatDefault:     simpletable.StyleDefault,
		formatCompact:     simpletable.StyleCompact,
		formatCompactLite: simpletable.StyleCompactLite,
		formatMarkdown:    simpletable.StyleMarkdown,
		formatRounded:     simpletable.StyleRounded,
		formatUnicode:     simpletable.StyleUnicode,
	}
)
