package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// NewReverseProxy creates an httputil.ReverseProxy that forwards requests to
// the given target base URL (e.g. "http://expenses-service:8081").
func NewReverseProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		panic("proxy: invalid upstream URL: " + err.Error())
	}
	return httputil.NewSingleHostReverseProxy(u)
}

// Handler wraps an httputil.ReverseProxy as a gin.HandlerFunc.
// Path params captured by Gin (e.g. /:id) are forwarded transparently because
// the ReverseProxy uses c.Request.URL, which Gin has already rewritten.
func Handler(p *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		p.ServeHTTP(c.Writer, c.Request)
	}
}

// ErrorHandler is an optional Director-level error handler for the proxy.
// It returns a plain JSON 502 when the upstream is unreachable.
func ErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"upstream service unavailable"}`))
	}
}
