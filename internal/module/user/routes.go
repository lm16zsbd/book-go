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

	router.Route("/books", func(sub chi.Router) {
		sub.Post("/{bookId}/favourite", httputil.Wrap(r.Handler.ToggleFavourite))
		sub.Post("/cancelFavourite", httputil.Wrap(r.Handler.BatchCancelFavourite))
		sub.Get("/{bookId}/bookmarks", httputil.Wrap(r.Handler.GetBookmarks))
		sub.Post("/{bookId}/bookmarks", httputil.Wrap(r.Handler.CreateBookmark))
		sub.Get("/{bookId}/annotations", httputil.Wrap(r.Handler.GetAnnotations))
		sub.Post("/{bookId}/annotations", httputil.Wrap(r.Handler.CreateAnnotation))
	})

	router.Delete("/bookmarks/{bookmarkId}", httputil.Wrap(r.Handler.DeleteBookmark))

	router.Route("/reader-config", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.GetReaderConfig))
		sub.Post("/", httputil.Wrap(r.Handler.SaveReaderConfig))
	})
}
