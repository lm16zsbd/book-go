package reader

import (
	"github.com/go-chi/chi/v5"

	"serica-go/internal/pkg/httputil"
)

const BasePath = "/v1/client/reader"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Route("/config", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.GetConfig))
		sub.Patch("/", httputil.Wrap(r.Handler.UpdateConfig))
	})

	router.Route("/notes", func(sub chi.Router) {
		sub.Post("/", httputil.Wrap(r.Handler.CreateAnnotation))
		sub.Patch("/{annotationId}", httputil.Wrap(r.Handler.UpdateAnnotation))
		sub.Delete("/{annotationId}", httputil.Wrap(r.Handler.DeleteAnnotation))
	})

	router.Get("/books/{id}/bookmarks", httputil.Wrap(r.Handler.GetBookmarks))
	router.Post("/bookmarks", httputil.Wrap(r.Handler.CreateBookmark))
	router.Get("/books/{id}/notes", httputil.Wrap(r.Handler.GetAnnotations))
}
