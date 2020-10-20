package progress

import (
	"math"
	"sync"

	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

type progress struct {
	p     *mpb.Progress
	b     *mpb.Bar
	m     sync.Mutex
	total int64
}

func (progress *progress) OnStepQueued() {
	progress.m.Lock()
	defer progress.m.Unlock()

	progress.total++
	progress.b.SetTotal(progress.total, false)
}

func (progress *progress) OnStepDone() {
	progress.m.Lock()
	defer progress.m.Unlock()

	progress.b.Increment()
}

func New() *progress {
	p := mpb.New(mpb.WithWidth(64))
	total := math.MaxInt64
	bar := p.AddBar(int64(total),
		mpb.AppendDecorators(decor.Percentage()),
	)

	return &progress{
		p: p,
		b: bar,
	}
}
