package async

import (
	"context"
	"sync"
)

func MultiRead[T any](ctx context.Context, ch ...<-chan T) <-chan T {

	out := make(chan T)

	go func() {
		defer close(out)
		wg := sync.WaitGroup{}
		for _, in := range ch {
			wg.Add(1)
			go func(wg *sync.WaitGroup, in <-chan T) {
				defer wg.Done()
				<-Pipe(ctx, out, in)
			}(&wg, in)
		}
		wg.Wait()
	}()

	return out
}
