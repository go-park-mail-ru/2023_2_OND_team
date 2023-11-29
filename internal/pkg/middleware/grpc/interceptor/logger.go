package interceptor

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"google.golang.org/grpc"
)

func Logger(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func(t time.Time) {
			log.Info("call rpc", logger.F{"handler", info.FullMethod}, logger.F{"time_execution", time.Since(t).Milliseconds()})
		}(time.Now())
		res, err := handler(ctx, req)
		return res, err
	}
}
