package workerpool

import (
	"fmt"
	"sync"
)

type work[I, O any] struct {
	t Task[O]
	r Result[O]
}

func newWork[I, O any](t Task[O]) *work[I, O] {
	result := Result[O]{
		r:   make(chan O, 1),
		err: make(chan error, 1),
	}

	return &work[I, O]{
		t: t,
		r: result,
	}
}

func (p *Pool[I, O]) Do(wg *sync.WaitGroup) {
	for {
		select {
		case <-p.ctx.Done():
			wg.Done()
			return
		case w := <-p.works:
			out, err := p.DoWithDefer(w.t)
			if err != nil {
				close(w.r.r)
				w.r.err <- err
				continue
			}
			close(w.r.err)
			w.r.r <- out
		}
	}
}

func (p *Pool[I, O]) DoWithDefer(t Task[O]) (_ O, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in task: %v", r)
		}
	}()

	return t(p.ctx)
}
