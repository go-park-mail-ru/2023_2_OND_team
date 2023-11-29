package main

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app/messenger"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func main() {
	log, err := logger.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	messenger.Run(context.Background(), log)
}
