package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	prefix          string
	totalHits       prometheus.Counter
	hitsDetail      *prometheus.CounterVec
	timeMeasurement *prometheus.HistogramVec
}

func (m metrics) AddRequest(method, path string, statusResponse int, executed time.Duration) {
	var labelStatus string
	switch {
	case statusResponse < 200:
		labelStatus = "100"
	case statusResponse < 300:
		labelStatus = "200"
	case statusResponse < 400:
		labelStatus = "300"
	case statusResponse < 500:
		labelStatus = "400"
	default:
		labelStatus = "500"
	}

	m.totalHits.Inc()
	m.hitsDetail.WithLabelValues(method, path, labelStatus).Inc()
	m.timeMeasurement.WithLabelValues(method, path, labelStatus).Observe(float64(executed.Milliseconds()))
}

func New(prefix string) metrics {
	return metrics{
		prefix: prefix,
		totalHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_total_hits",
			Help: "number of all requests",
		}),
		hitsDetail: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_hits_detail",
			Help: "the number of requests indicating its method, path, and response status",
		}, []string{"method", "path", "status"}),
		timeMeasurement: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_request_time_execution",
			Help:    "time of request execution with indication of its method, path and response status",
			Buckets: []float64{10, 100, 500, 1000, 5000},
		}, []string{"method", "path", "status"}),
	}
}

func (m metrics) Registry() (err error) {
	defer func() {
		if err != nil {
			prometheus.Unregister(m.totalHits)
			prometheus.Unregister(m.hitsDetail)
			prometheus.Unregister(m.timeMeasurement)
		}
	}()

	err = prometheus.Register(m.totalHits)
	if err != nil {
		return fmt.Errorf("registry metric total hits: %w", err)
	}

	err = prometheus.Register(m.hitsDetail)
	if err != nil {
		return fmt.Errorf("registry metric hits detail: %w", err)
	}

	err = prometheus.Register(m.timeMeasurement)
	if err != nil {
		return fmt.Errorf("registry metric time measurement: %w", err)
	}

	return nil
}
