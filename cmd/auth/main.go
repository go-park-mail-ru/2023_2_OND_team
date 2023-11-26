package main

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func main() {
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	auth.Run(context.Background(), log, configAuth)
}
