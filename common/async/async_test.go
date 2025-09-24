package async

import (
	"context"
	"testing"
	"time"
)

func TestPipe(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	t1 := time.NewTicker(time.Second)
	t2 := make(chan time.Time)
	defer close(t2)

	defer t.Log("done")

	done := Pipe(ctx, t2, t1.C)

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-done:
			break LOOP
		case v, ok := <-t2:
			if !ok {
				break LOOP
			}
			t.Log(v)
		}
	}

}
func TestMultiRead(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	t1 := time.NewTicker(time.Second)
	t2 := time.NewTicker(time.Millisecond * 400)

	defer t.Log("done")

	t3 := MultiRead(ctx, t1.C, t2.C)

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case v, ok := <-t3:
			if !ok {
				break LOOP
			}
			t.Log(v)
		}
	}

}
