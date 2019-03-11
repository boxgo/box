package loader

import "github.com/boxgo/box/minibox"

type (
	// Loader 配置加载器
	Loader interface {
		Load(minibox.MiniBox)
	}
)
