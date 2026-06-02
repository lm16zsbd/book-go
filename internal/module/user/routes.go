package user

import (
	"github.com/go-chi/chi/v5"

	"serica-go/internal/pkg/httputil"
)

const BasePath = "/v1/client"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Post("/login", httputil.Wrap(r.Handler.Login))

	router.Route("/user", func(sub chi.Router) {
		sub.Post("/register", httputil.Wrap(r.Handler.Register))
		sub.Get("/profile", httputil.Wrap(r.Handler.GetProfile))
		sub.Patch("/profile", httputil.Wrap(r.Handler.UpdateProfile))
		sub.Get("/favourites", httputil.Wrap(r.Handler.GetFavourites))
	})
}
