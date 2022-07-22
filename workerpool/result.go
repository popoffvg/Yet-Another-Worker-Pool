package workerpool

// Result represent calculating value. 
// It contains mathods for get value and errors, readiness check.
type Result[O any] struct {
	r        chan O
	err      chan error
	done     chan struct{}
	notEmpty bool
}

// If result is ready than return closed channel.
// Else return empty channel.
func (r *Result[O]) Done() <-chan struct{} {
	if r.isEmpty() {
		close(r.done)
	}

	return r.done
}

// Try to get value. If result is not ready than blocked.
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
