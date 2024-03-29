package ginprom

import "net/http"

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go#L112-L129
func computeApproximateRequestSize(r *http.Request) float64 {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return float64(s)
}
