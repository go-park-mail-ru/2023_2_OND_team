package main

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/app"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

const configFile = "configs/config.yml"

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
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	app.Run(log, configFile)
}
