package server

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrNotInitRouter = errors.New("there is no routing")

type Server struct {
	http.Server
	log *logger.Logger
	cfg *Config
}

func New(log *logger.Logger, cfg *Config) *Server {
	return &Server{
		Server: http.Server{
			Addr: cfg.Host + ":" + cfg.Port,
		},
		log: log,
		cfg: cfg,
	}
}

func (s *Server) Run(handler http.Handler) error {
	s.Handler = handler

	s.log.Info("server start")
	if s.cfg.https {
		return s.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	return s.ListenAndServe()
}
