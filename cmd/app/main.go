package main

import (
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/service"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func main() {
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	cfg, err := newConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	service := service.New(log, nil, nil, nil)
	server := server.New(log, *server.NewConfig(cfg))
	server.InitRouter(service)
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
