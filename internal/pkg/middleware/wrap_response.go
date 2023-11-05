package middleware

import (
	"io"
	"net/http"
)

type wrapResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
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

func (w *wrapResponseWriter) Write(data []byte) (written int, err error) {
	written, err = w.ResponseWriter.Write(data)
	w.written += written
	return
}

func (w *wrapResponseWriter) WriteString(data string) (written int, err error) {
	written, err = io.WriteString(w.ResponseWriter, data)
	w.written += written
	return
}
