package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_2_OND_team/docs"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
)

type Router struct {
	Mux *chi.Mux
}

func New() Router {
	return Router{chi.NewMux()}
}

func (r Router) RegisterRoute(handler *deliveryHTTP.HandlerHTTP) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://pinspire.online", "https://pinspire.online:1443"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders:   []string{"content-type"},
	})

	r.Mux.Use(auth.NewAuthMiddleware(nil, nil).Middleware, c.Handler)

	r.Mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/docs/*", httpSwagger.WrapHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Get("/login", handler.CheckLogin)
			r.Post("/login", handler.Login)
			r.Post("/signup", handler.Signup)
			r.Delete("/logout", handler.Logout)
		})

		r.Route("/pin", func(r chi.Router) {
			r.Get("/", handler.GetPins)
		})
	})
}
