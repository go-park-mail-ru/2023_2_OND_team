package app

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authProto "github.com/go-park-mail-ru/2023_2_OND_team/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	deliveryWS "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/websocket"
	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board/postgres"
	imgRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image"
	mesRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/message"
	pinRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	subRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/subscription/postgres"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var _timeoutForConnPG = 5 * time.Second

const uploadFiles = "upload/"

func Run(ctx context.Context, log *log.Logger, cfg ConfigFiles) {
	godotenv.Load()

	ctx, cancelCtxPG := context.WithTimeout(ctx, _timeoutForConnPG)
	defer cancelCtxPG()

	pool, err := NewPoolPG(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	imgCase := image.New(log, imgRepo.NewImageRepoFS(uploadFiles))
	messageCase := message.New(mesRepo.NewMessageRepo(pool))

	conn, err := grpc.Dial(cfg.AddrAuthServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer conn.Close()
	ac := auth.New(authProto.NewAuthClient(conn))

	handler := deliveryHTTP.New(log, deliveryHTTP.UsecaseHub{
		AuhtCase:         ac,
		UserCase:         user.New(log, imgCase, userRepo.NewUserRepoPG(pool)),
		PinCase:          pin.New(log, imgCase, pinRepo.NewPinRepoPG(pool)),
		BoardCase:        board.New(log, boardRepo.NewBoardRepoPG(pool), userRepo.NewUserRepoPG(pool), bluemonday.UGCPolicy()),
		SubscriptionCase: subscription.New(log, subRepo.NewSubscriptionRepoPG(pool), userRepo.NewUserRepoPG(pool)),
		MessageCase:      messageCase,
	})

	wsHandler := deliveryWS.New(log, messageCase,
		deliveryWS.SetOriginPatterns([]string{"pinspire.online", "pinspire.online:*"}))

	cfgServ, err := server.NewConfig(cfg.ServerConfigFile)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := server.New(log, cfgServ)
	router := router.New()
	router.RegisterRoute(handler, wsHandler, ac, log)

	if err := server.Run(router.Mux); err != nil {
		log.Error(err.Error())
		return
	}
}
