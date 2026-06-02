package home

import "github.com/go-chi/chi/v5"

const BasePath = "/v1/client"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Get("/home", r.Handler.Home)

	router.Route("/homepage", func(sub chi.Router) {
		sub.Get("/", r.Handler.Homepage)
		sub.Get("/section", r.Handler.HomepageSection)
	})
}
