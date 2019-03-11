package box

import (
	"context"

	"github.com/boxgo/box/config"
	"github.com/boxgo/box/minibox"
)

// WithContext 设置context到Box应用生命周期，默认是context.Background()
func WithContext(ctx context.Context) Option {
	return func(op *Options) {
		op.Context = ctx
	}
}

// WithBoxes 预置boxes
func WithBoxes(boxes []minibox.MiniBox) Option {
	return func(op *Options) {
		op.Boxes = boxes
	}
}

// WithConfig 设置应用配置信息管理器
func WithConfig(config config.Config) Option {
	return func(op *Options) {
		op.Config = config
	}
}

// WithConfigHook 设置全部配置信息获取Pre与Post钩子
func WithConfigHook(hook minibox.ConfigHook) Option {
	return func(op *Options) {
		op.ConfigHook = hook
	}
}

// WithServerHook 设置全局服务启动Pre与Post
func WithServerHook(hook minibox.ServerHook) Option {
	return func(op *Options) {
		op.ServerHook = hook
	}
}
