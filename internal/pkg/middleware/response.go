package middleware

import "net/http"

type wrapResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewWrapResponseWriter(w http.ResponseWriter) *wrapResponseWriter {
	return &wrapResponseWriter{
		ResponseWriter: w,
	}
}

func (w *wrapResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func SetResponseHeaders(headers map[string]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range headers {
				w.Header().Set(key, value)
			}
			next.ServeHTTP(w, r)
		})
	}
}
