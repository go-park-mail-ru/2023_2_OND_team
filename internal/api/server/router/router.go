package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_2_OND_team/docs"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/service"
)

type Router struct {
	Mux *chi.Mux
}

func New() Router {
	return Router{chi.NewMux()}
}

func (r Router) RegisterRoute(serv *service.Service) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://pinspire.online", "https://pinspire.online:1443"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders:   []string{"content-type"},
	})

	r.Mux.Use(c.Handler)

	r.Mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/docs/*", httpSwagger.WrapHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Get("/login", serv.CheckLogin)
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
