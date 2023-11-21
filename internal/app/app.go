package app

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board/postgres"
	imgRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image"
	pinRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	sessionRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session"
	subRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/subscription/postgres"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var (
	timeoutForConnPG    = 5 * time.Second
	timeoutForConnRedis = 5 * time.Second
)

const uploadFiles = "upload/"

func Run(ctx context.Context, log *log.Logger, cfg ConfigFiles) {
	godotenv.Load()

	ctx, cancelCtxPG := context.WithTimeout(ctx, timeoutForConnPG)
	defer cancelCtxPG()

	pool, err := NewPoolPG(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	ctx, cancelCtxRedis := context.WithTimeout(ctx, timeoutForConnRedis)
	defer cancelCtxRedis()

	redisCfg, err := NewConfig(cfg.RedisConfigFile)
	if err != nil {
		log.Error(err.Error())
		return
	}

	redisCl, err := NewRedisClient(ctx, redisCfg)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer redisCl.Close()

	sm := session.New(log, sessionRepo.NewSessionRepo(redisCl))
	imgCase := image.New(log, imgRepo.NewImageRepoFS(uploadFiles))
	userCase := user.New(log, imgCase, userRepo.NewUserRepoPG(pool))
	pinCase := pin.New(log, imgCase, pinRepo.NewPinRepoPG(pool))
	boardCase := board.New(log, boardRepo.NewBoardRepoPG(pool), userRepo.NewUserRepoPG(pool), bluemonday.UGCPolicy())
	subCase := subscription.New(log, subRepo.NewSubscriptionRepoPG(pool), userRepo.NewUserRepoPG(pool))

	handler := deliveryHTTP.New(log, sm, userCase, pinCase, boardCase, subCase)
	cfgServ, err := server.NewConfig(cfg.ServerConfigFile)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := server.New(log, cfgServ)
	router := router.New()
	router.RegisterRoute(handler, sm, log)

	if err := server.Run(router.Mux); err != nil {
		log.Error(err.Error())
		return
	}
}
