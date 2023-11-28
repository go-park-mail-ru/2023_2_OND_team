package monitoring

import (
	"net/http"
	"time"

	mw "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	AddRequest(method, path string, statusResponse int, executed time.Duration)
}

type StatusReceiver interface {
	StatusCode() int
}

func Monitoring(pathExporter string, metrics Metrics) func(http.Handler) http.Handler {
	instrumentMetricHandler := promhttp.Handler()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == pathExporter {
				instrumentMetricHandler.ServeHTTP(w, r)
				return
			}

			stat, ok := w.(StatusReceiver)
			if !ok {
				wrapResponse := mw.NewWrapResponseWriter(w)
				w = wrapResponse
				stat = wrapResponse
			}

			defer func(method, path string, t time.Time) {
				metrics.AddRequest(method, path, stat.StatusCode(), time.Since(t))
			}(r.Method, r.URL.Path, time.Now())

			next.ServeHTTP(w, r)
		})
	}
}
