package main

import "github.com/go-park-mail-ru/2023_2_OND_team/internal/app/auth"

var configAuth = auth.Config{
	Addr:            "localhost:8085",
	RedisFileConfig: "redis.conf",
}
