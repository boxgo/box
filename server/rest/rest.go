package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// Server supported by gin
	Server struct {
		Port int    `json:"port" desc:"Rest server serve port, default is 8080"`
		Addr string `json:"addr" desc:"Rest server listen addr, default is localhost"`
		Doc  bool   `json:"doc" desc:"Whether to generate an api document"`

		server *http.Server
		Engine *gin.Engine `json:"-"`
	}
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// Name box name
func (server *Server) Name() string {
	return "rest"
}

// ConfigWillLoad before config load
func (server *Server) ConfigWillLoad(ctx context.Context) {
	if server.Engine == nil {
		panic("Rest server must new by NewServer function.")
	}
}

// ConfigDidLoad after config load
func (server *Server) ConfigDidLoad(ctx context.Context) {
	if server.Port == 0 {
		server.Port = 8080
	}
}

// Serve box serve handler
func (server *Server) Serve(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", server.Addr, server.Port)

	server.server = &http.Server{
		Addr:    addr,
		Handler: server.Engine,
	}

	if err := server.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown box shutdown handler
func (server *Server) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

// NewServer new a rest server
func NewServer() *Server {
	server := &Server{
		Engine: gin.New(),
	}

	return server
}
