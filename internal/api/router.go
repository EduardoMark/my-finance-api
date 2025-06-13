package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Api) BindRoutes() {
	api.Router.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Route("/v1", func(r chi.Router) {
			api.Handler.User.RegisterRoutes(r)
		})

	})
}
