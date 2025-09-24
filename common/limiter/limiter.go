package limiter

import (
	"time"
)

// Limiter 令牌桶限速器
type Limiter struct {
	tokens      chan struct{} // 令牌桶
	duration    time.Duration // 颁发令牌周期
	assignCount int           // 每次颁发几个令牌，决定最小并发数
	concurrency int           // 令牌桶容量上限，决定最大并发数
	done        chan struct{} // <-Done()
}

func New(duration time.Duration, assignCount, concurrency int) *Limiter {
	return (&Limiter{
		duration:    duration,
		assignCount: assignCount,
		concurrency: concurrency,
	}).start(0)
}

func NewWithDelay(delay, duration time.Duration, assignCount, concurrency int) *Limiter {
	return (&Limiter{
		duration:    duration,
		assignCount: assignCount,
		concurrency: concurrency,
	}).start(delay)
}

func (l *Limiter) start(delay time.Duration) *Limiter {

	l.tokens = make(chan struct{}, l.concurrency)
	l.done = make(chan struct{})

	go func() {

		defer close(l.tokens)

		tick := time.NewTicker(l.duration)
		defer tick.Stop()

		for {
			select {
			case <-l.done:
				return
			case <-tick.C:
				l.assign()
			}
		}
	}()

	// 开始先分配一次，避免无效等待
	if delay > 0 {
		time.AfterFunc(delay, l.assign)
	} else {
		l.assign()
	}

	return l
}

// Stop 终止限速器
func (l *Limiter) Stop() *Limiter {
	close(l.done)
	return l
}

func (l *Limiter) assign() {
	for i := 0; i < l.assignCount; i++ {
		select {
		case <-l.done:
			return
		case l.tokens <- struct{}{}:
		default:
		}
	}
}

// Done 已停止生成新令牌
func (l *Limiter) Done() <-chan struct{} {
	return l.done
}

// Wait 阻塞等待
func (l *Limiter) Wait() {
	<-l.tokens
}

// Chan 获取令牌
func (l *Limiter) Chan() <-chan struct{} {
	return l.tokens
}
