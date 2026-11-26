package insight

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boxgo/box/pkg/server/ginserver"
)

func TestConfigRoutesRequireBasicAuth(t *testing.T) {
	accounts := map[string]string{"reader": "test-password"}
	for _, tt := range []struct {
		name        string
		serverAuth  map[string]string
		routeAuth   map[string]string
		validStatus int
	}{
		{name: "config credentials", routeAuth: accounts, validStatus: http.StatusOK},
		{name: "server credentials", serverAuth: accounts, validStatus: http.StatusOK},
		{name: "both", serverAuth: accounts, routeAuth: accounts, validStatus: http.StatusOK},
		{name: "both with different passwords", serverAuth: accounts, routeAuth: map[string]string{"reader": "other-password"}, validStatus: http.StatusUnauthorized},
		{name: "unconfigured", validStatus: http.StatusUnauthorized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ginserver.DefaultConfig("insight")
			cfg.BasicAuth = tt.serverAuth
			server := cfg.Build()
			registerConfigRoutes(server, len(cfg.BasicAuth) > 0, tt.routeAuth)

			for _, path := range []string{"/config", "/config?format=json", "/config?format=env", "/config/metric.basicAuth"} {
				for _, auth := range []struct {
					password string
					status   int
				}{
					{status: http.StatusUnauthorized},
					{password: "incorrect", status: http.StatusUnauthorized},
					{password: accounts["reader"], status: tt.validStatus},
				} {
					req := httptest.NewRequest(http.MethodGet, path, nil)
					if auth.password != "" {
						req.SetBasicAuth("reader", auth.password)
					}
					rec := httptest.NewRecorder()
					server.ServeHTTP(rec, req)
					if rec.Code != auth.status {
						t.Errorf("%s returned %d, want %d", path, rec.Code, auth.status)
					}
					if auth.status == http.StatusUnauthorized && rec.Header().Get("WWW-Authenticate") == "" {
						t.Errorf("%s did not challenge an unauthorized request", path)
					}
				}

				req := httptest.NewRequest(http.MethodGet, path, nil)
				req.SetBasicAuth("reader", accounts["reader"])
				header := strings.Replace(req.Header.Get("Authorization"), "Basic ", "basic ", 1)
				req.Header.Set("Authorization", header)
				rec := httptest.NewRecorder()
				server.ServeHTTP(rec, req)
				if rec.Code != tt.validStatus {
					t.Errorf("%s with lowercase basic scheme returned %d, want %d", path, rec.Code, tt.validStatus)
				}
			}
		})
	}
}
