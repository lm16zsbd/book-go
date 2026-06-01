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
	router.Get("/homepage", r.Handler.Homepage)
	router.Get("/homepage/section", r.Handler.HomepageSection)
}
