package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	api "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/api"
	handler "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/delivery/grpc/v1"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/repository"
	roll "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/service"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc"
)

const (
	maxConnDB        = 500
	schemaDB         = "public"
	timeoutForConnPG = 5 * time.Second
)

func NewPoolPG(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsnPG())

	cfg.MaxConns = maxConnDB
	cfg.ConnConfig.RuntimeParams["search_path"] = schemaDB

	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new postgres pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping new pool: %w", err)
	}
	return pool, err
}

func dsnPG() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
}

func main() {
	godotenv.Load()
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	ctx, cancelCtxPG := context.WithTimeout(context.Background(), timeoutForConnPG)
	defer cancelCtxPG()

	pool, err := NewPoolPG(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer pool.Close()

	rollService := roll.New(log, repo.NewRollRepoPG(pool), bluemonday.UGCPolicy())

	lis, err := net.Listen("tcp", ":8100")
	if err != nil {
		log.Fatal("can't listen port")
	}
	grpcServer := handler.NewServerGRPC(log, rollService)

	server := grpc.NewServer()
	api.RegisterRollServiceServer(server, grpcServer)
	fmt.Println("START SERVER")
	server.Serve(lis)
}

/*
lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	session.RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting server at :8081")
	server.Serve(lis)
*/
