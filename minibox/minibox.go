package minibox

import (
	"context"
)

type (
	// MiniBox a mini box
	// 迷你盒子是一个微型应用，Duck Typing优雅的解决了已注册迷你盒子对功能的自由支持与扩展。
	// * 支持从配置中心获取自己以及外部的配置信息
	// * 支持配置信息获取过程的Pre与Post钩子
	// * 支持优雅的启动与停止迷你盒子的服务
	// * 支持迷你盒子启停过程的Pre与Post钩子
	MiniBox interface {
		// Name 盒子的名称
		// 根据名称加载指定路径的配置信息，支持以点（.）作为访问层级分隔符，支持以空字符串作为访问根层级标识符
		// name: ""     				加载配置信息中根目录下的配置信息到当前迷你盒子
		// name: "test" 				加载配置信息中test路径下的配置信息到当前迷你盒子
		// name: "test.example" 加载配置信息中test下example的配置信息到当前迷你盒子
		Name() string
	}

	// MiniBoxExt 一般情况下迷你盒子仅可以获取自己Name下的配置信息，但是如果实现了此接口即可以获取外部迷你盒子的配置信息
	MiniBoxExt interface {
		Exts() []MiniBox
	}

	// ConfigHook 迷你盒子配置信息获取过程的Pre与Post钩子
	ConfigHook interface {
		ConfigWillLoad(ctx context.Context)
		ConfigDidLoad(ctx context.Context)
	}

	// Server 迷你盒子启动与停止过程处理器
	Server interface {
		Serve(ctx context.Context) error
		Shutdown(ctx context.Context) error
	}

	// ServerHook 迷你盒子启停过程的Pre与Post钩子
	ServerHook interface {
		ServerWillReady(ctx context.Context)
		ServerDidReady(ctx context.Context)
		ServerWillClose(ctx context.Context)
		ServerDidClose(ctx context.Context)
	}
)
