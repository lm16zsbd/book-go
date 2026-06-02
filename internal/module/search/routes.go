package search

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
	router.Route("/search", func(sub chi.Router) {
		sub.Get("/history", httputil.Wrap(r.Handler.GetHistory))
		sub.Delete("/history", httputil.Wrap(r.Handler.ClearHistory))
		sub.Get("/hot", httputil.Wrap(r.Handler.GetHot))
	})

	router.Get("/hot-word", httputil.Wrap(r.Handler.GetHotWords))
	router.Get("/search-word", httputil.Wrap(r.Handler.GetRecommendWords))
}
