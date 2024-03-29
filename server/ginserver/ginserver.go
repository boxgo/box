package ginserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	GinServer struct {
		cfg    *Config
		engine *gin.Engine
		server *http.Server
	}
)

var (
	default404Body = []byte(http.StatusText(http.StatusNotFound))
	default405Body = []byte(http.StatusText(http.StatusMethodNotAllowed))
)

func newGinServer(cfg *Config) *GinServer {
	gin.SetMode(gin.ReleaseMode) // init mode is release.

	engine := gin.New()
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	gin.SetMode(cfg.Mode)

	if cfg.BasicAuth != nil && len(cfg.BasicAuth) != 0 {
		engine.Use(gin.BasicAuth(cfg.BasicAuth))
	}

	engine.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusNotFound, "text/plain", default404Body)
	})
	engine.NoMethod(func(c *gin.Context) {
		c.Data(http.StatusMethodNotAllowed, "text/plain", default405Body)
	})

	return &GinServer{
		cfg:    cfg,
		engine: engine,
		server: server,
	}
}

func (server *GinServer) Name() string {
	return "gin-server"
}

func (server *GinServer) Serve(ctx context.Context) error {
	return server.Run()
}

func (server *GinServer) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

func (server *GinServer) Run() error {
	if err := server.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (server *GinServer) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.Use(middleware...)
}

func (server *GinServer) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.Any(relativePath, handlers...)
}

func (server *GinServer) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.DELETE(relativePath, handlers...)
}

func (server *GinServer) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.GET(relativePath, handlers...)
}

func (server *GinServer) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.HEAD(relativePath, handlers...)
}

func (server *GinServer) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.OPTIONS(relativePath, handlers...)
}

func (server *GinServer) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.PATCH(relativePath, handlers...)
}

func (server *GinServer) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.POST(relativePath, handlers...)
}

func (server *GinServer) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.engine.PUT(relativePath, handlers...)
}

func (server *GinServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	server.engine.ServeHTTP(w, req)
}

func (server *GinServer) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return server.engine.Group(relativePath, handlers...)
}

func (server *GinServer) RoutesInfo() gin.RoutesInfo {
	return server.engine.Routes()
}
