package insight

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/server/ginserver"
	"github.com/gin-gonic/gin"
)

var (
	Default = ginserver.StdConfig("insight").Build()
)

func init() {
	Get("/", func(ctx *gin.Context) {
		ctx.Data(200, gin.MIMEHTML, []byte(html()))
	})

	Get("/config", func(ctx *gin.Context) {
		switch ctx.Query("format") {
		case "json":
			ctx.JSON(200, config.Fields())
		case "table":
			ctx.Data(200, gin.MIMEPlain, bytes.NewBufferString(config.Fields().Table()).Bytes())
		case "env":
			ctx.Data(200, gin.MIMEPlain, bytes.NewBufferString(config.Fields().Env()).Bytes())
		default:
			ctx.Data(200, gin.MIMEPlain, bytes.NewBufferString(config.Fields().Table()).Bytes())
		}
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

func html() string {
	routes := Default.RoutesInfo()
	paths := make([]string, len(routes))

	for idx, info := range routes {
		if info.Path == "/" {
			continue
		}

		paths[idx] = fmt.Sprintf(`<li><a href="%s">%s</a></li>`, info.Path, info.Path)
	}

	sort.Strings(paths)

	return fmt.Sprintf(`
<html>
<head>
<style>
    body {
        width: 35em;
        margin: 0 auto;
    }
</style>
</head>
<body>
<h1>Welcome to %s %s insight</h1>
<ul>%s</ul>
</body>
</html>`, config.ServiceName(), config.ServiceVersion(), strings.Join(paths, ""))
}
