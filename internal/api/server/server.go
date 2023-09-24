package server

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/api/server/router"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/service"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrNotInitRouter = errors.New("there is no routing")

type Server struct {
	http.Server
	router router.Router
	log    *logger.Logger
}

func New(log *logger.Logger, cfg Config) *Server {
	return &Server{
		Server: http.Server{
			Addr: cfg.Host + ":" + cfg.Port,
		},
		router: router.New(),
		log:    log,
	}
}

func (s *Server) Run() error {
	if !s.router.IsInit() {
		return ErrNotInitRouter
	}
	s.Handler = s.router.Mux
	s.log.Info("server start")
	return s.ListenAndServe()
}

func (s *Server) InitRouter(serv *service.Service) {
	s.router.InitRoute(serv)
}
