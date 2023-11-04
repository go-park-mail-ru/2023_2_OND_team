package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type Middleware func(next http.Handler) http.Handler

type ctxKeyRequestID string

const RequestIDKey ctxKeyRequestID = "RequestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.NewRandom()
		if err != nil {
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestIDKey, id.String())))
	})
}

func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, ok := r.Context().Value(RequestIDKey).(string)
			if !ok {
				log.Warn("for middleware to work with logging, enable middleware with request id assignment")
				next.ServeHTTP(w, r)
				return
			}

			wrapResponse := NewWrapResponseWriter(w)

			log.Info("request", logger.F{"request-id", requestID},
				logger.F{"method", r.Method}, logger.F{"path", r.URL.Path},
				logger.F{"content-type", r.Header.Get("Content-Type")},
				logger.F{"content-length", r.Header.Get("Content-Length")},
				logger.F{"address", r.RemoteAddr})

			defer func(t time.Time) {
				log.Info("response", logger.F{"request-id", requestID},
					logger.F{"status", strconv.FormatInt(int64(wrapResponse.statusCode), 10)},
					logger.F{"processing_time", strconv.FormatInt(int64(time.Since(t).Milliseconds()), 10) + "ms"},
					logger.F{"content-type", w.Header().Get("Content-Type")},
					logger.F{"content-length", w.Header().Get("Content-Length")})
			}(time.Now())
			next.ServeHTTP(w, r)
		})
	}
}
