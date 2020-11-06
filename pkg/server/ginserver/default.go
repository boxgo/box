package ginserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() error {
	return Default.Run()
}

func Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return Default.Use(middleware...)
}

func Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.Any(relativePath, handlers...)
}

func DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.DELETE(relativePath, handlers...)
}

func GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.GET(relativePath, handlers...)
}

func HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.HEAD(relativePath, handlers...)
}

func OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.OPTIONS(relativePath, handlers...)
}

func PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.PATCH(relativePath, handlers...)
}

func POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.POST(relativePath, handlers...)
}

func PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return Default.PUT(relativePath, handlers...)
}

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	Default.ServeHTTP(w, req)
}
