package insight

import (
	"bytes"
	"net/http"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/server/ginserver"
	"github.com/gin-gonic/gin"
)

type (
	routes []route
	route  struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	}
)

var (
	Default = ginserver.StdConfig("insight").Build()
)

func init() {
	Get("/", func(ctx *gin.Context) {
		rs := make(routes, len(Default.RoutesInfo()))

		for idx, info := range Default.RoutesInfo() {
			rs[idx] = route{
				Method: info.Method,
				Path:   info.Path,
			}
		}

		ctx.JSON(200, rs)
	})

	Get("/config", func(ctx *gin.Context) {
		ctx.JSON(200, config.Fields())
	})

	Get("/config/table", func(ctx *gin.Context) {
		ctx.Data(200, gin.MIMEPlain, bytes.NewBufferString(config.Fields().Table()).Bytes())
	})
}

func Any(relativePath string, handlers ...gin.HandlerFunc) {
	Default.Any(relativePath, handlers...)
}

func AnyH(relativePath string, handler http.Handler) {
	Default.Any(relativePath, gin.WrapH(handler))
}

func AnyF(relativePath string, handler http.HandlerFunc) {
	Default.Any(relativePath, gin.WrapF(handler))
}

func Get(relativePath string, handlers ...gin.HandlerFunc) {
	Default.GET(relativePath, handlers...)
}

func GetH(relativePath string, handler http.Handler) {
	Default.GET(relativePath, gin.WrapH(handler))
}

func GetF(relativePath string, handler http.HandlerFunc) {
	Default.GET(relativePath, gin.WrapF(handler))
}

func Post(relativePath string, handlers ...gin.HandlerFunc) {
	Default.POST(relativePath, handlers...)
}

func PostH(relativePath string, handler http.Handler) {
	Default.POST(relativePath, gin.WrapH(handler))
}

func PostF(relativePath string, handler http.HandlerFunc) {
	Default.POST(relativePath, gin.WrapF(handler))
}
