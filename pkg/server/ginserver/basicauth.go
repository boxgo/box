package ginserver

import (
	"net/http"

	"github.com/boxgo/box/internal/basicauth"
	"github.com/gin-gonic/gin"
)

// BasicAuth protects a Gin route or group. Without accounts, it rejects every request.
func BasicAuth(accounts map[string]string, realm string) gin.HandlerFunc {
	verify := basicauth.NewVerifier(accounts)
	challenge := basicauth.Challenge(realm)
	return func(ctx *gin.Context) {
		username, ok := verify(ctx.Request)
		if !ok {
			ctx.Header("WWW-Authenticate", challenge)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(gin.AuthUserKey, username)
	}
}
