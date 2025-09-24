package callback

import (
	"context"
	"github.com/jack0829/letsgo/queue"
	"io"
	"net/http"
	"time"
)

type (
	option struct {
		c []callbackOption
		q []queue.Option[*Task]
	}
	callbackOption func(c *Callback)
	Option         func(o *option)
)

func (o *option) New(ctx context.Context) *Callback {

	q := queue.New[*Task](ctx, o.q...)
	c := &Callback{
		queue: q,
	}

	for _, op := range o.c {
		op(c)
	}

	return c
}

func WithClient(client *http.Client) Option {
	return func(o *option) {
		o.c = append(o.c, func(c *Callback) {
			c.client = client
		})
	}
}

func WithSaver(s Saver) Option {
	return func(o *option) {
		o.q = append(
			o.q,
			queue.WithSaveHandler(s.Save),
			queue.WithData(s.Load()...),
		)
	}
}

func WithRetry(
	max int, // <=0：不重试；1：重试1次；2：重试2次；...
	delay time.Duration, // 重试间隔，至少 5 秒
) Option {
	return func(o *option) {
		o.c = append(o.c, func(c *Callback) {
			if max > 0 {
				c.retryMax = max
				c.retryDelay = time.Second * 5
				if delay > c.retryDelay {
					c.retryDelay = delay
				}
			}
		})
	}
}

func WithID(id string) TaskOption {
	return func(t *Task) {
		t.ID = id
	}
}

func WithBody(r io.Reader) TaskOption {
	return func(t *Task) {
		t.body = r
	}
}

func WithBodyString(s string) TaskOption {
	return func(t *Task) {
		t.Body = s
	}
}
