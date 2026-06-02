package search

import "github.com/go-chi/chi/v5"

const BasePath = "/v1/client"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Route("/search", func(sub chi.Router) {
		sub.Get("/history", r.Handler.GetHistory)
		sub.Delete("/history", r.Handler.ClearHistory)
		sub.Get("/hot", r.Handler.GetHot)
	})

	router.Get("/hot-word", r.Handler.GetHotWords)
	router.Get("/search-word", r.Handler.GetRecommendWords)
}
