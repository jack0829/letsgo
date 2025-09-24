package callback

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	FH "github.com/jack0829/letsgo/http"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"
)

type testRequest struct {
	I    uint64 `json:"i"`
	Time string `json:"time"`
}

func testTaskStream(
	ctx context.Context,
	d time.Duration,
) <-chan *Task {

	ch := make(chan *Task)
	go func() {
		defer close(ch)
		tick := time.NewTicker(d)
		defer tick.Stop()

		var i uint64
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
			}

			i++
			// body := fmt.Sprintf(`{"i":%d,"time":"%s"}`, i, time.Now().Format(time.TimeOnly))
			body := bytes.NewBuffer(nil)
			jsoniter.NewEncoder(body).Encode(testRequest{
				I:    i,
				Time: time.Now().Format(time.TimeOnly),
			})
			select {
			case <-ctx.Done():
			case ch <- NewTask("http://127.0.0.1:8080/callback", WithBody(body)):
			}
		}
	}()

	return ch
}

func testHttpServer(
	ctx context.Context,
) {

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		var r testRequest
		if err := jsoniter.NewDecoder(req.Body).Decode(&r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.I%3 == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsoniter.NewEncoder(w).Encode(r)
		return
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		defer s.Shutdown(context.TODO())
		<-ctx.Done()
	}()

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	fmt.Println("HTTP-Server shutdown")
}

func TestCallback(t *testing.T) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go testHttpServer(ctx)

	cb := New(
		ctx,
		WithRetry(2, time.Second*5),
		WithClient(&http.Client{
			Transport: (&FH.Transport{}).Debug(os.Stdout),
		}),
	)

	go func() {
		t.Log("start listen")
		cb.Listen(ctx, 3)
		t.Log("stop listen")
	}()

	for v := range testTaskStream(ctx, time.Second*10) {
		if err := cb.Do(v); err != nil {
			t.Errorf("cb.Do err: %v", err)
		}
	}

	t.Log("结束")

}
