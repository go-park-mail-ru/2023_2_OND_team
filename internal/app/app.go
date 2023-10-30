package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
	pinRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	sessionRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func Run(ctx context.Context, log *log.Logger, configFile string) {
<<<<<<< HEAD
	pool, err := pgxpool.New(ctx, "postgres://ond_team:love@localhost:5432/pinspire?search_path=pinspire")
=======
	ctxApp, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctxApp, "postgres://ond_team:love@localhost:5432/pinspire?search_path=pinspire")
>>>>>>> dev2
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	err = pool.Ping(ctxApp)
	if err != nil {
		log.Error(err.Error())
		return
	}

	redisCl := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "love",
	})

	status := redisCl.Ping(ctxApp)
	if status.Err() != nil {
		log.Error(err.Error())
		return
	}

	sm := session.New(log, sessionRepo.NewSessionRepo(redisCl))
	userCase := user.New(log, userRepo.NewUserRepoPG(pool))
	pinCase := pin.New(log, pinRepo.NewPinRepoPG(pool))
	boardCase := board.New(log, boardRepo.NewBoardRepoPG(pool))

	handler := deliveryHTTP.New(log, sm, userCase, pinCase, boardCase)
	cfgServ, err := server.NewConfig(configFile)
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
