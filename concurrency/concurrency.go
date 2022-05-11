package concurrency

import (
	"context"
	"runtime"
	"sync"
)

type Concurrency struct {
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  func()
	errOnce sync.Once

	err error

	qCh chan func() error
}

func WithContext(ctx context.Context, limit int) (cc *Concurrency, c context.Context) {
	if limit == 0 {
		limit = runtime.NumCPU()
	}
	c, cancel := context.WithCancel(ctx)
	cc = &Concurrency{ctx: c, cancel: cancel, qCh: make(chan func() error, 0)}
	cc.wg.Add(limit)
	for i := 0; i < limit; i++ {
		go cc.routine()
	}
	return
}

func (g *Concurrency) routine() {
	defer g.wg.Done()
	for {
		select {
		case <-g.ctx.Done():
			return
		case f, ok := <-g.qCh:
			if !ok {
				return
			}
			if f == nil {
				continue
			}
			if err := f(); err != nil {
				g.errOnce.Do(func() {
					g.cancel()
					g.err = err
				})
			}
		}
	}
}

func (g *Concurrency) SubmitDone() {
	close(g.qCh)
}

func (g *Concurrency) Go(f func() error) {
	select {
	case <-g.ctx.Done():
	case g.qCh <- f:
	}
}

func (g *Concurrency) Wait() (err error) {
	close(g.qCh)
	g.wg.Wait()

	if g.err != nil {
		return g.err
	}
	return g.ctx.Err()
}
