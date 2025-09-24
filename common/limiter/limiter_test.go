package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {

	l := New(time.Second*1, 2, 5)
	wg := &sync.WaitGroup{}

	test := func(id int, after time.Duration) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-l.Done():
				// l.Wait()
				t.Logf("Closed\t%d\t%s", id, time.Now().Format("15:04:05"))
			case a := <-time.After(after):
				t.Logf("Waiting\t%d\t%s", id, a.Format("15:04:05"))
				l.Wait()
				t.Logf("Waited\t%d\t%s", id, time.Now().Format("15:04:05"))
			}
		}()
	}

	for i := 1; i <= 10; i++ {
		test(i, time.Second*time.Duration(i+2))
	}

	time.AfterFunc(time.Second*5, func() {
		l.Stop()
	})

	wg.Wait()
	t.Log("Done")
}
