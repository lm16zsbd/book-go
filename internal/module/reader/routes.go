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
	router.Get("/config", r.Handler.GetConfig)
	router.Patch("/config", r.Handler.UpdateConfig)

	router.Get("/books/{id}/bookmarks", r.Handler.GetBookmarks)
	router.Post("/bookmarks", r.Handler.CreateBookmark)

	router.Get("/books/{id}/notes", r.Handler.GetAnnotations)
	router.Post("/notes", r.Handler.CreateAnnotation)
	router.Patch("/notes/{annotationId}", r.Handler.UpdateAnnotation)
	router.Delete("/notes/{annotationId}", r.Handler.DeleteAnnotation)
}
