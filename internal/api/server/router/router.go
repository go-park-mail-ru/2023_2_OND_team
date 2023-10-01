package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/service"
)

type Router struct {
	Mux *chi.Mux
}

func New() Router {
	return Router{chi.NewMux()}
}

func (r Router) InitRoute(serv *service.Service) {
	r.Mux.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", serv.Login)
			r.Post("/signup", serv.Signup)
			r.Delete("/logout", serv.Logout)
		})

		r.Route("/pin", func(r chi.Router) {
			r.Get("/", serv.GetPins)
		})
	})
}

func (r Router) IsInit() bool {
	return r.Mux != nil
}
