package async

import (
	"context"
)

func Pipe[T any](ctx context.Context, w chan<- T, r <-chan T) <-chan struct{} {

	done := make(chan struct{})

	go func() {

		defer close(done)

		var (
			v  T
			ok bool
		)

	LOOP:
		for {

			select {
			case <-ctx.Done():
				break LOOP
			case v, ok = <-r:
				if !ok {
					break LOOP
				}
			}

			select {
			case w <- v:
			case <-ctx.Done():
				break LOOP
			}
		}

	}()

	return done
}
