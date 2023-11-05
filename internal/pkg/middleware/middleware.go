package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type Middleware func(next http.Handler) http.Handler

type ctxKeyRequestID string

const RequestIDKey ctxKeyRequestID = "RequestID"

func RequestID(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := uuid.NewRandom()
			if err != nil {
				log.Sugar().Errorf("middleware requestID: %s", err.Error())
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestIDKey, id.String())))
		})
	}
}

func Logger(logMW *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logMW
			if requestID, ok := r.Context().Value(RequestIDKey).(string); ok {
				log = log.WithField("request-id", requestID)
			} else {
				log.Warn("for middleware to work with logging, enable middleware with request id assignment")
			}

			wrapResponse := NewWrapResponseWriter(w)

			log.InfoMap("request", logger.M{
				"method":         r.Method,
				"path":           r.URL.Path,
				"content_type":   r.Header.Get("Content-Type"),
				"content_length": r.ContentLength,
				"address":        r.RemoteAddr,
			})
			defer func(t time.Time) {
				log.InfoMap("response", logger.M{
					"status":             wrapResponse.statusCode,
					"processing_time_ms": time.Since(t).Milliseconds(),
					"content_type":       w.Header().Get("Content-Type"),
					"content_length":     w.Header().Get("Content-Length"),
					"written":            wrapResponse.written,
				})
			}(time.Now())
			next.ServeHTTP(wrapResponse, r.WithContext(context.WithValue(r.Context(), logger.KeyLogger, log)))
		})
	}
}
