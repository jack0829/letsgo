package bootstrap

import (
	"context"
	"fmt"
	"sync"
)

type Register struct {
	handlers        []Handler
	destroyHandlers []func()
	executed        *sync.Once
	destroyed       *sync.Once
	wg              *sync.WaitGroup
}

// Handler 初始化句柄，返回一个销毁方法和初始化时产生的错误
type Handler func(ctx context.Context) (destroy func(), err error)

func New() *Register {
	return &Register{
		handlers:        make([]Handler, 0),
		destroyHandlers: make([]func(), 0),
		executed:        &sync.Once{},
		destroyed:       &sync.Once{},
		wg:              &sync.WaitGroup{},
	}
}

func (r *Register) Add(fn Handler) {
	r.handlers = append(r.handlers, fn)
}

// StartWithContext 执行启动项
func (r *Register) StartWithContext(
	ctx context.Context,
) (err error) {

	r.executed.Do(func() {

		var destroy func()
		for _, fn := range r.handlers {

			destroy, err = fn(ctx)
			if err != nil {
				return
			}

			if destroy != nil {
				r.destroyHandlers = append(r.destroyHandlers, destroy)
				r.wg.Add(1)
			}
		}
	})

	return
}

// End 执行启动项释放动作
func (r *Register) End() {

	r.destroyed.Do(func() {

		defer func() {
			if err := recover(); err != nil {
				fmt.Println("bootstrap deregister error", err)
			}
		}()

		if handlers := r.destroyHandlers; len(handlers) > 0 {
			for _, fn := range handlers {
				fn()
				r.wg.Done()
			}
		}
	})

	r.wg.Wait()
}

// Start 默认 Context 的启动
func (r *Register) Start() error {
	return r.StartWithContext(context.TODO())
}
