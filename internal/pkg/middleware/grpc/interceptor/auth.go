package interceptor

import (
	"context"
	"strconv"

	messMS "github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/messenger/delivery/grpc"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md := metadata.ValueFromIncomingContext(ctx, messMS.AuthenticatedMetadataKey)
		if len(md) != 1 {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}
		userID, err := strconv.ParseInt(md[0], 10, 64)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}
		return handler(context.WithValue(ctx, auth.KeyCurrentUserID, int(userID)), req)
	}
}
