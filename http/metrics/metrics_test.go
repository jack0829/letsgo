package metrics

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	randDuration := func(min, max time.Duration) time.Duration {
		return min + time.Duration(rand.Int63n(int64(max-min)))
	}

	m := New("test", WithLabels(map[string]string{
		"instance": "My Macbook",
		"foo":      "bar",
	}))

	r := gin.New()

	r.GET("/metrics", gin.WrapH(m.Exporter()))

	r.GET("/hello", m.Gin, func(g *gin.Context) {
		d := randDuration(time.Second, time.Second*3)
		time.Sleep(d)
		g.String(http.StatusOK, "hello world ! %s", d)
		// g.Next()
	})

	srv := &http.Server{Handler: r}
	go func() {

		l, err := net.Listen("tcp", ":8080")
		if err != nil {
			t.Error(err)
			return
		}
		defer l.Close()

		if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
		t.Log("srv.Shutdown")
	}()

	httpGet := func(url string) {
		resp, err := http.Get(url)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		b := bytes.NewBuffer(nil)
		if _, err := b.ReadFrom(resp.Body); err != nil {
			t.Error(err)
			return
		}
		t.Logf("%s %s", resp.Status, b.String())
	}

	go func() {
		tick := time.NewTicker(time.Second * 5)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 300):
				httpGet("http://127.0.0.1:8080/hello")
			case <-tick.C:
				httpGet("http://127.0.0.1:8080/metrics")
			}
		}
	}()

	<-ctx.Done()
	httpGet("http://127.0.0.1:8080/metrics")
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Error(err)
		return
	}
	t.Log("Done")
}
