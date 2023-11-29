package messenger

import (
	"context"
	"net"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app"
	messMS "github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/messenger/delivery/grpc"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/messenger/usecase/message"
	grpcMetrics "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/metrics/grpc"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/grpc/interceptor"
	mesRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

const _timeoutForConnPG = 5 * time.Second

func Run(ctx context.Context, log *logger.Logger) {
	godotenv.Load()

	metrics := grpcMetrics.New("messenger")
	if err := metrics.Registry(); err != nil {
		log.Error(err.Error())
		return
	}

	ctx, cancelCtxPG := context.WithTimeout(ctx, _timeoutForConnPG)
	defer cancelCtxPG()

	pool, err := app.NewPoolPG(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	messageCase := message.New(mesRepo.NewMessageRepo(pool))

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.Monitoring(metrics, "localhost:8096"),
		interceptor.Logger(log),
		interceptor.Auth(),
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
