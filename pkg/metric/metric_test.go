package metric

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandlerBasicAuth(t *testing.T) {
	protected := (&Config{BasicAuth: map[string]string{
		"prometheus": "scrape:password",
		"backup":     "other-password",
	}}).Build().Handler()
	public := (&Config{}).Build().Handler()

	tests := []struct {
		name               string
		handler            http.Handler
		username, password string
		authorization      string
		status             int
	}{
		{name: "unconfigured", handler: public, status: http.StatusOK},
		{name: "missing credentials", handler: protected, status: http.StatusUnauthorized},
		{name: "invalid header", handler: protected, authorization: "Basic !!!", status: http.StatusUnauthorized},
		{name: "wrong username", handler: protected, username: "unknown", password: "scrape:password", status: http.StatusUnauthorized},
		{name: "wrong password", handler: protected, username: "prometheus", password: "wrong", status: http.StatusUnauthorized},
		{name: "mixed accounts", handler: protected, username: "prometheus", password: "other-password", status: http.StatusUnauthorized},
		{name: "valid credentials", handler: protected, username: "prometheus", password: "scrape:password", status: http.StatusOK},
		{name: "second account", handler: protected, username: "backup", password: "other-password", status: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
			if tt.username != "" {
				req.SetBasicAuth(tt.username, tt.password)
			} else if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			rec := httptest.NewRecorder()
			tt.handler.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
			challenge := rec.Header().Get("WWW-Authenticate")
			hasMetrics := strings.Contains(rec.Body.String(), "go_goroutines")
			if tt.status == http.StatusUnauthorized {
				if challenge != `Basic realm="metrics"` || hasMetrics {
					t.Fatalf("challenge = %q, metrics exposed = %t", challenge, hasMetrics)
				}
			} else if challenge != "" || !hasMetrics {
				t.Fatalf("challenge = %q, metrics returned = %t", challenge, hasMetrics)
			}
		})
	}
}

func TestPushBasicAuth(t *testing.T) {
	for _, tt := range []struct {
		name string
		auth PushBasicAuth
	}{
		{name: "authenticated", auth: PushBasicAuth{Username: "pusher", Password: "push:secret"}},
		{name: "anonymous"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			requests := make(chan *http.Request, 1)
			gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				select {
				case requests <- r:
				default:
				}
				w.WriteHeader(http.StatusAccepted)
			}))
			defer gateway.Close()

			cfg := DefaultConfig("")
			cfg.PushEnabled = true
			cfg.PushTargetURL = gateway.URL
			cfg.PushInterval = 10 * time.Millisecond
			cfg.PushBasicAuth = tt.auth
			metric := cfg.Build()
			if err := metric.Serve(context.Background()); err != nil {
				t.Fatal(err)
			}
			defer metric.Shutdown(context.Background())

			select {
			case req := <-requests:
				username, password, ok := req.BasicAuth()
				if ok != (tt.auth != (PushBasicAuth{})) || username != tt.auth.Username || password != tt.auth.Password {
					t.Fatal("Pushgateway received unexpected Basic Auth credentials")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Pushgateway did not receive a metrics request")
			}
		})
	}
}
