package interceptor

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type Metrics interface {
	AddRequest(handler string, ok bool, executed time.Duration)
}

func Monitoring(m Metrics, addres string) grpc.UnaryServerInterceptor {
	serv := http.Server{
		Addr:    addres,
		Handler: promhttp.Handler(),
	}
	go func() { serv.ListenAndServe() }()
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var ok = true
		defer func(t time.Time) {
			m.AddRequest(info.FullMethod, ok, time.Since(t))
		}(time.Now())
		res, err := handler(ctx, req)
		if err != nil {
			ok = false
		}
		return res, err
	}
}
