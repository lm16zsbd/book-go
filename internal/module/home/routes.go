package home

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
	router.Get("/home", httputil.Wrap(r.Handler.Home))

	router.Route("/homepage", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.Homepage))
		sub.Get("/section", httputil.Wrap(r.Handler.HomepageSection))
	})
}
