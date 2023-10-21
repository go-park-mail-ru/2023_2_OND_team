package app

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func Run(log *log.Logger, configFile string) {
	db, err := ramrepo.OpenDB("RamRepository")
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer db.Close()

	sm := session.New(log, ramrepo.NewRamSessionRepo(db))
	userCase := user.New(log, ramrepo.NewRamUserRepo(db))
	pinCase := pin.New(log, ramrepo.NewRamPinRepo(db))

	handler := deliveryHTTP.New(log, sm, userCase, pinCase)
	cfgServ, err := server.NewConfig(configFile)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := server.New(log, cfgServ)
	router := router.New()
	router.RegisterRoute(handler, sm, log)

	if err := server.Run(router.Mux); err != nil {
		log.Error(err.Error())
		return
	}
}
