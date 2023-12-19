package main

import (
	"os"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app"
)

var configFiles = app.ConfigFiles{
	ServerConfigFile: "configs/config.yml",
	AddrAuthServer:   os.Getenv("AUTH_SERVICE_HOST") + ":" + os.Getenv("AUTH_SERVICE_PORT"), // "localhost:8085",
}
