package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SetRequestTimeout(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if clientTimeout := extractTimeout(r); clientTimeout != 0 {
				timeout = time.Duration(clientTimeout)
			}
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTimeout(r *http.Request) int {
	var keepAlive string
	if keepAlive = r.Header.Get("Keep-Alive"); keepAlive == "" {
		return *new(int)
	}

	if options := strings.Split(keepAlive, " "); len(options) > 0 && len(options) <= 2 {
		if timeoutOpt := options[0]; strings.Contains(timeoutOpt, "timeout") {
			if timeout, err := strconv.ParseInt(strings.Split(timeoutOpt, "=")[1], 10, 64); err == nil {
				return int(timeout)
			}
		}
	}
	return *new(int)
}
