package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2023_2_OND_team/docs"
	deliveryHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/security"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type Router struct {
	Mux *chi.Mux
}

func New() Router {
	return Router{chi.NewMux()}
}

func (r Router) RegisterRoute(handler *deliveryHTTP.HandlerHTTP, sm session.SessionManager, log *logger.Logger) {
	cfgCSRF := security.DefaultCSRFConfig()
	cfgCSRF.PathToGet = "/api/v1/csrf"

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://pinspire.online", "https://pinspire.online:1443",
			"https://pinspire.online:1444", "https://pinspire.online:1445", "https://pinspire.online:1446"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowCredentials: true,
		AllowedHeaders:   []string{"content-type", cfgCSRF.Header},
		ExposedHeaders:   []string{cfgCSRF.HeaderSet},
	})

	r.Mux.Use(c.Handler, auth.NewAuthMiddleware(sm).ContextWithUserID) //security.CSRF(cfgCSRF)

	r.Mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/docs/*", httpSwagger.WrapHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handler.Login)
			r.Post("/signup", handler.Signup)
			r.Group(func(r chi.Router) {
				r.Use(auth.RequireAuth)

				r.Get("/login", handler.CheckLogin)
				r.Delete("/logout", handler.Logout)
			})
		})

		r.With(auth.RequireAuth).Route("/profile", func(r chi.Router) {
			r.Get("/info", handler.GetProfileInfo)
			r.Put("/edit", handler.ProfileEditInfo)
			r.Put("/avatar", handler.ProfileEditAvatar)
		})

		r.Route("/pin", func(r chi.Router) {
			r.Get("/", handler.GetPins)
			r.Group(func(r chi.Router) {
				r.Use(auth.RequireAuth)

				r.Post("/create", handler.CreateNewPin)
				r.Post("/like/{pinID:\\d+}", handler.SetLikePin)
				r.Put("/edit/{pinID:\\d+}", handler.EditPin)
				r.Delete("/like/{pinID:\\d+}", handler.DeleteLikePin)
				r.Delete("/delete/{pinID:\\d+}", handler.DeletePin)
			})
		})

		r.Route("/board", func(r chi.Router) {
			r.Route("/get", func(r chi.Router) {
				r.Get("/user/{username}", handler.GetUserBoards)
				r.Get("/{boardID:\\d+}", handler.GetCertainBoard)
			})
			r.Group(func(r chi.Router) {
				r.Use(auth.RequireAuth)
				r.Post("/create", handler.CreateNewBoard)
				r.Put("/update/{boardID:\\d+}", handler.UpdateBoardInfo)
				r.Delete("/delete/{boardID:\\d+}", handler.DeleteBoard)
			})
		})
	})
}
