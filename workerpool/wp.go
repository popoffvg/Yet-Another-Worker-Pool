package workerpool

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrRedirectOut = errors.New("wrong source of result")
)

type Task[O any] func(context.Context) (O, error)

type Worker[I, O any] func(context.Context, I) (O, error)

type Option[I, O any] func(*Pool[I, O])

func WithWorker[I, O any](w Worker[I, O]) Option[I, O] {
	return func(p *Pool[I, O]) {
		p.workerTemplate = w
	}
}

func WithCtx[I, O any](ctx context.Context) Option[I, O] {
	return func(p *Pool[I, O]) {
		ctx, cancel := context.WithCancel(ctx)
		p.cancel = cancel
		p.ctx = ctx
	}
}

func RedirectOutput[I, O any]() Option[I, O] {
	return func(p *Pool[I, O]) {
		p.out = make(chan Result[O], 1)
	}
}

type Pool[I, O any] struct {
	works          chan work[I, O]
	workerTemplate Worker[I, O]
	ctx            context.Context

	cancel   context.CancelFunc
	cancelWG sync.WaitGroup

	out chan Result[O]
}

func New[I, O any](workers int, opts ...Option[I, O]) *Pool[I, O] {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool[I, O]{
		works:  make(chan work[I, O], 1),
		ctx:    ctx,
		cancel: cancel,
	}

	for _, option := range opts {
		option(p)
	}

	for i := 0; i < workers; i++ {
		p.cancelWG.Add(1)
		go p.Do(&p.cancelWG)
	}

	return p
}

func (p *Pool[I, O]) Pub(input I) Result[O] {
	if p.workerTemplate == nil {
		panic("not set worker template")
	}
	var task Task[O] = func(ctx context.Context) (O, error) {
		return p.workerTemplate(ctx, input)
	}
	if p.out != nil {
		return Result[O]{}
	}
	return p.Run(task)
}

func (p *Pool[I, O]) Run(t Task[O]) Result[O] {
	work := newWork[I](t)
	p.works <- *work
	if p.out != nil {
		return Result[O]{}
	}
	return work.r
}

func (p *Pool[I, O]) Close() error {
	select {
	case <-p.ctx.Done():
		return nil
	default:
	}
	p.cancel()
	p.cancelWG.Wait()
	close(p.out)

	return nil
}

func (p *Pool[I, O]) Stream() <-chan Result[O] {
	if p.out == nil {
		mock := make(chan Result[O])
		close(mock)
		return mock
	}

	return p.out
}