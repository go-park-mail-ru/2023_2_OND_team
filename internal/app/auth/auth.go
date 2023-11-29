package auth

import (
	"context"
	"net"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	authProto "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app"
	authMS "github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/auth/delivery/grpc"
	sessRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var (
	_timeoutForConnPG    = 5 * time.Second
	_timeoutForConnRedis = 5 * time.Second
)

func Run(ctx context.Context, log *logger.Logger, cfg Config) {
	godotenv.Load()

	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer l.Close()

	ctxPG, cancelCtxPG := context.WithTimeout(ctx, _timeoutForConnPG)
	defer cancelCtxPG()

	pool, err := app.NewPoolPG(ctxPG)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	ctxRedis, cancelCtxRedis := context.WithTimeout(ctx, _timeoutForConnRedis)
	defer cancelCtxRedis()

	redisCfg, err := app.NewConfig(cfg.RedisFileConfig)
	if err != nil {
		log.Error(err.Error())
		return
	}

	redisCl, err := app.NewRedisClient(ctxRedis, redisCfg)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer redisCl.Close()

	sm := session.New(log, sessRepo.NewSessionRepo(redisCl))
	u := user.New(log, nil, userRepo.NewUserRepoPG(pool))

	s := grpc.NewServer()
	authProto.RegisterAuthServer(s, authMS.New(log, sm, u))

	log.Info("service auht start", logger.F{"addr", cfg.Addr})
	if err = s.Serve(l); err != nil {
		log.Error(err.Error())
		return
	}
}