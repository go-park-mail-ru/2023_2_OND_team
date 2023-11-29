package messenger

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app"
	messMS "github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/messenger/delivery/grpc"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	mesRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

const _timeoutForConnPG = 5 * time.Second

func Run(ctx context.Context, log *logger.Logger) {
	godotenv.Load()

	ctx, cancelCtxPG := context.WithTimeout(ctx, _timeoutForConnPG)
	defer cancelCtxPG()

	pool, err := app.NewPoolPG(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	messageCase := message.New(mesRepo.NewMessageRepo(pool))

	server := grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			f := metadata.ValueFromIncomingContext(ctx, string(auth.KeyCurrentUserID))
			if len(f) != 1 {
				return nil, status.Error(codes.Unauthenticated, "unauthenticated")
			}
			userID, err := strconv.ParseInt(f[0], 10, 64)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, "unauthenticated")
			}
			return handler(context.WithValue(ctx, auth.KeyCurrentUserID, int(userID)), req)
		},
	))
	messenger.RegisterMessengerServer(server, messMS.New(log, messageCase))

	l, err := net.Listen("tcp", "localhost:8095")
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("server messenger start", logger.F{"addr", "localhost:8095"})
	if err := server.Serve(l); err != nil {
		log.Error(err.Error())
		return
	}
}
