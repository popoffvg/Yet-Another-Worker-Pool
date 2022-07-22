package workerpool

type Result[O any] struct {
	r        chan O
	err      chan error
	done     chan struct{}
	notEmpty bool
}

func (r *Result[O]) Done() <-chan struct{} {
	if r.isEmpty() {
		close(r.done)
	}

	return r.done
}

func (r *Result[O]) Get() (O, error) {
	if r.isEmpty() {
		return getZero[O](), ErrRedirectOut
	}

	err := <-r.err
	if err != nil {
		return getZero[O](), err
	}

	return <-r.r, nil
}

func (r *Result[O]) set(out O, err error) {
	r.notEmpty = true
	close(r.done)
	r.err <- err
	r.r <- out
}

func (r *Result[O]) isEmpty() bool {
	return !r.notEmpty
}

func getZero[T any]() T {
	var result T
	return result
}
