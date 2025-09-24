package bootstrap

import (
	"context"
)

var (
	defaultRegistry *Register
)

func init() {
	defaultRegistry = New()
}

// Add 添加启动项
func Add(fn Handler) {
	defaultRegistry.Add(fn)
}

// Start 执行启动项
func Start() error {
	return defaultRegistry.Start()
}

// StartWithContext 执行启动项
func StartWithContext(ctx context.Context) error {
	return defaultRegistry.StartWithContext(ctx)
}

// End 执行启动项释放动作
func End() {
	defaultRegistry.End()
}
