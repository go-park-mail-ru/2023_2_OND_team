package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2023_2_OND_team/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app"
	authMS "github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/auth/delivery/grpc"
	imgRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image"
	sessionRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := net.Listen("tcp", "localhost:8085")
	if err != nil {
		log.Error(err.Error())
		return
	}

	s := grpc.NewServer()
	godotenv.Load()

	ctx := context.Background()

	pool, err := app.NewPoolPG(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	ctx, cancelCtxRedis := context.WithTimeout(ctx, time.Second)
	defer cancelCtxRedis()

	redisCfg, err := app.NewConfig("redis.conf")
	if err != nil {
		log.Error(err.Error())
		return
	}

	redisCl, err := app.NewRedisClient(ctx, redisCfg)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer redisCl.Close()

	sm := session.New(log, sessionRepo.NewSessionRepo(redisCl))
	imgCase := image.New(log, imgRepo.NewImageRepoFS("upload/"))
	u := user.New(log, imgCase, userRepo.NewUserRepoPG(pool))
	auth.RegisterAuthServer(s, authMS.New(log, sm, u))
	if err = s.Serve(l); err != nil {
		log.Error(err.Error())
		return
	}
}
