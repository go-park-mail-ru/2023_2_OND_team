package app

import (
	"context"
	"encoding/base64"
	"os"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	vision "cloud.google.com/go/vision/v2/apiv1"
	authProto "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	deliveryWS "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/websocket"
	notify "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/notification"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/metrics"
	commentNotify "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/notification/comment"
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
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/realtime/chat"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/realtime/notification"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/search"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	validate "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var (
	_timeoutForConnPG     = 5 * time.Second
	timeoutCloudVisionAPI = 10 * time.Second
)

const uploadFiles = "upload/"

func Run(ctx context.Context, log *log.Logger, cfg ConfigFiles) {
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

	connMessMS, err := grpc.Dial(os.Getenv("MESSENGER_SERVICE_HOST")+":"+os.Getenv("MESSENGER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer connMessMS.Close()

	connRealtime, err := grpc.Dial(os.Getenv("REALTIME_SERVICE_HOST")+":"+os.Getenv("REALTIME_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer connRealtime.Close()

	rtClient := rt.NewRealTimeClient(connRealtime)

	commentRepository := commentRepo.NewCommentRepoPG(pool)

	visionCtx, cancel := context.WithTimeout(ctx, timeoutCloudVisionAPI)
	defer cancel()

	token, err := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Error(err.Error())
		return
	}
	visionClient, err := vision.NewImageAnnotatorClient(visionCtx, option.WithCredentialsJSON(token))
	if err != nil {
		log.Error(err.Error())
		return
	}

	profanityCensor := goaway.NewProfanityDetector().WithCustomDictionary(
		append(goaway.DefaultProfanities, validate.GetLabels()...),
		goaway.DefaultFalsePositives,
		goaway.DefaultFalseNegatives,
	)

	imgCase := image.New(log, imgRepo.NewImageRepoFS(uploadFiles), image.NewFilter(visionClient, validate.NewCensor(profanityCensor)))
	messageCase := message.New(log, messenger.NewMessengerClient(connMessMS), chat.New(realtime.NewRealTimeChatClient(rtClient), log))
	pinCase := pin.New(log, imgCase, pinRepo.NewPinRepoPG(pool))

	notifyBuilder, err := notify.NewWithType(notify.NotifyComment)
	if err != nil {
		log.Error(err.Error())
		return
	}

	notifyCase := notification.New(realtime.NewRealTimeNotificationClient(rtClient), log,
		notification.Register(commentNotify.NewCommentNotify(notifyBuilder, comment.New(commentRepository, pinCase, nil), pinCase)))

	conn, err := grpc.Dial(cfg.AddrAuthServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer conn.Close()
	ac := auth.New(authProto.NewAuthClient(conn))

	handler := deliveryHTTP.New(log, deliveryHTTP.NewConverterHTTP(validate.NewSanitizerXSS(bluemonday.UGCPolicy()), validate.NewCensor(profanityCensor)), deliveryHTTP.UsecaseHub{
		AuhtCase:         ac,
		UserCase:         user.New(log, imgCase, userRepo.NewUserRepoPG(pool)),
		PinCase:          pinCase,
		BoardCase:        board.New(log, boardRepo.NewBoardRepoPG(pool), userRepo.NewUserRepoPG(pool)),
		SubscriptionCase: subscription.New(log, subRepo.NewSubscriptionRepoPG(pool), userRepo.NewUserRepoPG(pool)),
		SearchCase:       search.New(log, searchRepo.NewSearchRepoPG(pool)),
		MessageCase:      messageCase,
		CommentCase:      comment.New(commentRepo.NewCommentRepoPG(pool), pinCase, notifyCase),
	})

	wsHandler := deliveryWS.New(log, messageCase, notifyCase,
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
