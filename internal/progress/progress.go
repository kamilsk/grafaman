package progress

import (
	"math"
	"sync/atomic"

	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

const defaultTerminalWidth = 80

type progress struct {
	p     *mpb.Progress
	b     *mpb.Bar
	total int64
}

func (progress *progress) OnStepQueued() {
	newTotal := atomic.AddInt64(&progress.total, 1)
	progress.b.SetTotal(newTotal, false)
}

func (progress *progress) OnStepDone() {
	progress.b.Increment()
}

func New() *progress {
	p := mpb.New(mpb.WithWidth(defaultTerminalWidth))
	bar := p.AddBar(math.MaxInt64, mpb.AppendDecorators(decor.Percentage()))

	return &progress{
		p: p,
		b: bar,
	}
}
