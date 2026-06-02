package reader

import "github.com/go-chi/chi/v5"

const BasePath = "/v1/client/reader"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Route("/config", func(sub chi.Router) {
		sub.Get("/", r.Handler.GetConfig)
		sub.Patch("/", r.Handler.UpdateConfig)
	})

	router.Route("/notes", func(sub chi.Router) {
		sub.Post("/", r.Handler.CreateAnnotation)
		sub.Patch("/{annotationId}", r.Handler.UpdateAnnotation)
		sub.Delete("/{annotationId}", r.Handler.DeleteAnnotation)
	})

	router.Get("/books/{id}/bookmarks", r.Handler.GetBookmarks)
	router.Post("/bookmarks", r.Handler.CreateBookmark)
	router.Get("/books/{id}/notes", r.Handler.GetAnnotations)
}
