package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type (
	SaveHandler[T any] func(val []T)
	Option[T any]      func(q *Queue[T])
)

type Queue[T any] struct {
	ctx     context.Context
	mutex   sync.Mutex
	waiting []T
	reading chan T
	saver   SaveHandler[T]
}

func New[T any](
	ctx context.Context,
	ops ...Option[T],
) *Queue[T] {

	reading := make(chan T)

	q := &Queue[T]{
		ctx:     ctx,
		reading: reading,
	}

	for _, op := range ops {
		op(q)
	}

	go func() {
		<-ctx.Done()
		close(reading)
		q.save(q.read())
	}()

	return q
}

func (q *Queue[T]) read() (list []T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.waiting) > 0 {
		list = q.waiting[:]
		q.waiting = nil
	}
	return
}

func (q *Queue[T]) write(v T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.waiting = append(q.waiting, v)
}

func (q *Queue[T]) Read(
	ctx context.Context,
	interval time.Duration,
) <-chan T {

	ch := make(chan T)

	go func() {

		defer close(ch)

		for {

			for _, v := range q.read() {
				ch <- v
			}

			select {
			case <-q.ctx.Done():
				return
			case <-ctx.Done():
				return
			case v := <-q.reading:
				ch <- v
			case <-time.After(interval):
			}

		}
	}()

	return ch
}

func (q *Queue[T]) Write(v T) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("queue: %s", r)
		}
	}()

	select {
	case <-q.ctx.Done():
	case q.reading <- v:
		// fmt.Println("chan<-", v)
	default:
		q.write(v)
		// fmt.Println("slice<-", v)
	}

	return
}

func (q *Queue[T]) save(data []T) {
	if q.saver != nil {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("queue save err: %v\n", r)
			}
		}()
		q.saver(data)
	}
}

func WithData[T any](data ...T) Option[T] {
	return func(q *Queue[T]) {
		if len(data) > 0 {
			q.mutex.Lock()
			q.waiting = data[:]
			q.mutex.Unlock()
		}
	}
}

func WithSaveHandler[T any](handler SaveHandler[T]) Option[T] {
	return func(q *Queue[T]) {
		q.saver = handler
	}
}
