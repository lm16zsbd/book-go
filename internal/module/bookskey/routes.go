package bookskey

import "github.com/go-chi/chi/v5"

const BasePath = "/v1/client"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Get("/booksKey", r.Handler.GetKey)
}
