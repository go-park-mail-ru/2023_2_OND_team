package messenger

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	prefix          string
	totalHits       prometheus.Counter
	hitsDetailOk    *prometheus.CounterVec
	hitsDetailErr   *prometheus.CounterVec
	timeMeasurement *prometheus.HistogramVec
}

func (m metrics) AddRequest(handler string, ok bool, executed time.Duration) {
	m.totalHits.Inc()
	if ok {
		m.hitsDetailOk.WithLabelValues(handler).Inc()
	} else {
		m.hitsDetailErr.WithLabelValues(handler).Inc()
	}
	m.timeMeasurement.WithLabelValues(handler).Observe(float64(executed.Milliseconds()))
}

func New(prefix string) metrics {
	return metrics{
		prefix: prefix,
		totalHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_total_hits",
			Help: "number of all requests",
		}),
		hitsDetailOk: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_hits_detail_ok",
			Help: "the number of requests indicating handler with normal status",
		}, []string{"handler"}),
		hitsDetailErr: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_hits_detail_err",
			Help: "the number of requests indicating its handler with error status",
		}, []string{"handler"}),
		timeMeasurement: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_request_time_execution",
			Help:    "time of request execution with indication of its handler",
			Buckets: []float64{10, 100, 500, 1000, 5000},
		}, []string{"handler"}),
	}
}

func (m metrics) Registry() (err error) {
	defer func() {
		if err != nil {
			prometheus.Unregister(m.totalHits)
			prometheus.Unregister(m.hitsDetailOk)
			prometheus.Unregister(m.timeMeasurement)
			prometheus.Unregister(m.hitsDetailErr)
		}
	}()

	err = prometheus.Register(m.totalHits)
	if err != nil {
		return fmt.Errorf("registry metric total hits: %w", err)
	}

	err = prometheus.Register(m.hitsDetailOk)
	if err != nil {
		return fmt.Errorf("registry metric hits detail ok: %w", err)
	}

	err = prometheus.Register(m.hitsDetailErr)
	if err != nil {
		return fmt.Errorf("registry metric hits detail error: %w", err)
	}

	err = prometheus.Register(m.timeMeasurement)
	if err != nil {
		return fmt.Errorf("registry metric time measurement: %w", err)
	}

	return nil
}
