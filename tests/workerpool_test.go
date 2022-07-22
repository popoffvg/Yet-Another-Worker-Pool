package workerpool

import (
	"context"
	"sync"
	"testing"

	. "github.com/popoffvg/wp/workerpool"
	"github.com/stretchr/testify/assert"
)

func TestWorkerPoolRun(t *testing.T) {
	wp := New[int, int](3)
	result := make([]Result[int], 0, 1_000)
	for i := 0; i < 1_000; i++ {
		number := i + 1
		result = append(result, wp.Run(func(ctx context.Context) (int, error) {
			return number, nil
		}))
	}

	for _, v := range result {
		v, err := v.Get()
		assert.Nil(t, err)
		assert.NotEqual(t, 0, v)
	}
}
func TestWorkerPoolPub(t *testing.T) {
	wp := New(3,
		WithWorker(func(ctx context.Context, i int) (int, error) {
			return i + 1, nil
		}),
	)
	result := make([]Result[int], 0, 1_000)
	for i := 0; i < 1_000; i++ {
		number := i + 1
		result = append(result, wp.Pub(number))
	}

	for _, v := range result {
		v, err := v.Get()
		assert.Nil(t, err)
		assert.NotEqual(t, 0, v)
	}
}

func TestWorkerPoolCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wp := New(3,
		WithCtx[int, int](ctx),
	)
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		wp.Run(func(ctx context.Context) (int, error) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return 0, nil
				default:
				}
			}
		})
	}
	cancel()
	wp.Close()
	wg.Wait()
}
