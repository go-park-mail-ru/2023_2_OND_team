package main

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/service"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/user"
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

	db, err := ramrepo.OpenDB("RamRepository")
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer db.Close()

	sm := session.New(log, ramrepo.NewRamSessionRepo(db))
	userCase := user.New(log, ramrepo.NewRamUserRepo(db))
	pinCase := pin.New(log, ramrepo.NewRamPinRepo(db))

	service := service.New(log, sm, userCase, pinCase)
	cfgServ, err := server.NewConfig(configFile)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := server.New(log, cfgServ)
	server.InitRouter(service)
	if err := server.Run(); err != nil {
		log.Error(err.Error())
		return
	}
}
