package app

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authProto "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	deliveryWS "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/websocket"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/metrics"
	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board/postgres"
	commentRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/comment"
	imgRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image"
	pinRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	searchRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/search/postgres"
	subRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/subscription/postgres"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/search"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var _timeoutForConnPG = 5 * time.Second

const uploadFiles = "upload/"

func Run(ctx context.Context, log *log.Logger, cfg ConfigFiles) {
	godotenv.Load()

	metrics := metrics.New("pinspire")
	err := metrics.Registry()
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx, cancelCtxPG := context.WithTimeout(ctx, _timeoutForConnPG)
	defer cancelCtxPG()

	pool, err := NewPoolPG(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	connMessMS, err := grpc.Dial("localhost:8095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer connMessMS.Close()

	imgCase := image.New(log, imgRepo.NewImageRepoFS(uploadFiles))
	messageCase := message.New(messenger.NewMessengerClient(connMessMS))
	pinCase := pin.New(log, imgCase, pinRepo.NewPinRepoPG(pool))

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
		PinCase:          pinCase,
		BoardCase:        board.New(log, boardRepo.NewBoardRepoPG(pool), userRepo.NewUserRepoPG(pool), bluemonday.UGCPolicy()),
		SubscriptionCase: subscription.New(log, subRepo.NewSubscriptionRepoPG(pool), userRepo.NewUserRepoPG(pool), bluemonday.UGCPolicy()),
		SearchCase:       search.New(log, searchRepo.NewSearchRepoPG(pool), bluemonday.UGCPolicy()),
		MessageCase:      messageCase,
		CommentCase:      comment.New(commentRepo.NewCommentRepoPG(pool), pinCase),
	})

	connRealtime, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer connRealtime.Close()

	wsHandler := deliveryWS.New(log, messageCase, rt.NewRealTimeClient(connRealtime),
		deliveryWS.SetOriginPatterns([]string{"pinspire.online", "pinspire.online:*"}))

	cfgServ, err := server.NewConfig(cfg.ServerConfigFile)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := server.New(log, cfgServ)
	router := router.New()
	router.RegisterRoute(handler, wsHandler, ac, metrics, log)

	if err := server.Run(router.Mux); err != nil {
		log.Error(err.Error())
		return
	}
}
