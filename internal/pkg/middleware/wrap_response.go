package middleware

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"
)

var ErrUnimplementedMethod = errors.New("unimplemented method")

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

func (w *wrapResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, ErrUnimplementedMethod
}
