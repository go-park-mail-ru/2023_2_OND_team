package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/joho/godotenv"
)

var (
	logOutput      = flag.String("log", "stdout", "file paths to write logging output to")
	logErrorOutput = flag.String("logerror", "stderr", "path to write internal logger errors to.")
)

//	@title			Pinspire API
//	@version		1.0
//	@description	API for Pinspire project
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	godotenv.Load()
	flag.Parse()
	ctxBase, cancel := context.WithCancel(context.Background())
	defer cancel()

	log, err := logger.New(
		logger.RFC3339FormatTime(),
		logger.SetOutputPaths(*logOutput),
		logger.SetErrorOutputPaths(*logErrorOutput),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	configFiles := app.ConfigFiles{
		ServerConfigFile: "configs/config.yml",
		AddrAuthServer:   os.Getenv("AUTH_SERVICE_HOST") + ":" + os.Getenv("AUTH_SERVICE_PORT"), // "localhost:8085",
	}
	app.Run(ctxBase, log, configFiles)
}
