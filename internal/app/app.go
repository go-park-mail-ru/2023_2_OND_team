package app

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	deliveryWS "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/websocket"
	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board/postgres"
	imgRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image"
	mesCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/message"
	pinRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	sessionRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
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

	handler := deliveryHTTP.New(log, deliveryHTTP.UsecaseHub{
		UserCase:    user.New(log, imgCase, userRepo.NewUserRepoPG(pool)),
		PinCase:     pin.New(log, imgCase, pinRepo.NewPinRepoPG(pool)),
		BoardCase:   board.New(log, boardRepo.NewBoardRepoPG(pool), userRepo.NewUserRepoPG(pool), bluemonday.UGCPolicy()),
		MessageCase: message.New(mesCase.NewMessageRepo(pool)),
		SM:          sm,
	})

	wsHandler := deliveryWS.New(log)

	cfgServ, err := server.NewConfig(cfg.ServerConfigFile)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := server.New(log, cfgServ)
	router := router.New()
	router.RegisterRoute(handler, wsHandler, sm, log)

	if err := server.Run(router.Mux); err != nil {
		log.Error(err.Error())
		return
	}
}
