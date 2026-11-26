package insight

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"slices"
	"sort"
	"strings"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/config/field"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/server/ginserver"
	"github.com/gin-gonic/gin"
)

var (
	serverConfig = ginserver.StdConfig("insight")
	Default      = serverConfig.Build()
)

type configAccess struct {
	BasicAuth map[string]string `config:"configBasicAuth" desc:"Basic Auth accounts for Insight config routes"`
}

func (*configAccess) Path() string { return "gin.insight" }

func init() {
	Get("/", func(ctx *gin.Context) {
		ctx.Data(200, gin.MIMEHTML, []byte(html()))
	})

	Get("/buildInfo", func(ctx *gin.Context) {
		buildInfoStr := "Unable to determine version information."

		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			buildInfoStr = buildInfo.String()
		}

		ctx.Data(200, gin.MIMEPlain, []byte(buildInfoStr))
	})

	access := &configAccess{}
	if err := config.Scan(access); err != nil {
		logger.Panicf("Insight config access load error: %s", err)
	}
	registerConfigRoutes(Default, len(serverConfig.BasicAuth) > 0, access.BasicAuth)
}

func registerConfigRoutes(server *ginserver.GinServer, hasServerAuth bool, accounts map[string]string) {
	routes := server.Group("/config")
	if len(accounts) > 0 || !hasServerAuth {
		routes.Use(ginserver.BasicAuth(accounts, "insight config"))
	}

	routes.GET("", func(ctx *gin.Context) {
		fields := slices.Clone(*config.Fields())
		switch ctx.Query("format") {
		case "json":
			type configField struct {
				field.Field
				Value any `json:"value"`
			}
			values := make([]configField, 0, len(fields))
			for _, f := range fields {
				values = append(values, configField{Field: f, Value: f.Value.Interface()})
			}
			ctx.JSON(200, values)
		case "env":
			ctx.String(200, "%s", fields.Env())
		default:
			ctx.String(200, "%s", fields.Table())
		}
	})

	routes.GET("/:key", func(ctx *gin.Context) {
		ctx.Data(200, gin.MIMEPlain, config.Get(ctx.Param("key")).Bytes())
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

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	Default.ServeHTTP(w, req)
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
