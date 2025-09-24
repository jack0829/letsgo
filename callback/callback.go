package callback

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jack0829/letsgo/common/async"
	"github.com/jack0829/letsgo/queue"
	"io"
	"net/http"
	"strings"
	"time"
)

type (
	Callback struct {
		ctx        context.Context
		listen     chan *Task
		retryMax   int
		retryDelay time.Duration
		client     *http.Client
		queue      *queue.Queue[*Task]
	}
)

func New(
	ctx context.Context,
	ops ...Option,
) *Callback {

	var o option
	for _, op := range ops {
		op(&o)
	}

	return o.New(ctx)
}

func (c *Callback) stream(
	ctx context.Context,
	bufSize int,
) <-chan *Task {

	c.ctx = ctx
	c.listen = make(chan *Task, bufSize)
	go func() {
		<-ctx.Done()
		close(c.listen)
		c.listen = nil
	}()

	ch := make(chan *Task, bufSize)
	go func() {
		defer close(ch)
		for t := range c.queue.Read(ctx, time.Second*15) {

			// fmt.Printf("new retry: %s\t%d\n", t.Body, t.Status.Tries)
			// 延迟处理
			if d := time.Since(t.Status.ReqTime) + c.retryDelay; d > 0 {
				// fmt.Printf("delay: %s\n", d.String())
				select {
				case <-time.After(d):
				case <-ctx.Done():
				}
			}

			// fmt.Printf("retry: %s\n", t.Body)
			// 放进 ch
			select {
			case <-ctx.Done():
			case ch <- t:
			}
		}
	}()

	return async.MultiRead(ctx, c.listen, ch)
}

func (c *Callback) Listen(
	ctx context.Context,
	bufSize int,
) {

	for t := range c.stream(ctx, bufSize) {

		if t == nil {
			continue
		}

		if err := c.do(t); err != nil {
			if t.Status.Tries < c.retryMax {
				t.Status.ReqTime = time.Now()
				t.Status.Tries++
				t.Status.Error = err.Error()
				c.queue.Write(t)
			}
		}
	}
}

func (c *Callback) Do(t *Task) error {

	if c.listen == nil {
		return fmt.Errorf("callback 尚未启动监听")
	}

	select {
	case <-c.ctx.Done():
	case c.listen <- t:
	}

	return nil
}

func (c *Callback) do(t *Task) error {

	if t == nil {
		return nil
	}

	var body io.Reader
	if t.body != nil {
		buf := bytes.NewBuffer(nil)
		body = io.TeeReader(t.body, buf)
		defer func() {
			t.Body = buf.String()
		}()
	} else {
		body = strings.NewReader(t.Body)
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, t.URL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c().Do(req)
	if err != nil {
		return err
	}
	defer closeResp(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d | %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (c *Callback) c() *http.Client {
	if c.client != nil {
		return c.client
	}
	return http.DefaultClient
}

func closeResp(r *http.Response) {
	if r.Body != nil {
		r.Body.Close()
	}
}
