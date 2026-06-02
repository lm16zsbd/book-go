package user

import "github.com/go-chi/chi/v5"

const BasePath = "/v1/client"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Post("/login", r.Handler.Login)

	router.Route("/user", func(sub chi.Router) {
		sub.Post("/register", r.Handler.Register)
		sub.Get("/profile", r.Handler.GetProfile)
		sub.Patch("/profile", r.Handler.UpdateProfile)
		sub.Get("/favourites", r.Handler.GetFavourites)
	})

	router.Route("/books", func(sub chi.Router) {
		sub.Post("/{bookId}/favourite", r.Handler.ToggleFavourite)
		sub.Post("/cancelFavourite", r.Handler.BatchCancelFavourite)
		sub.Get("/{bookId}/bookmarks", r.Handler.GetBookmarks)
		sub.Post("/{bookId}/bookmarks", r.Handler.CreateBookmark)
		sub.Get("/{bookId}/annotations", r.Handler.GetAnnotations)
		sub.Post("/{bookId}/annotations", r.Handler.CreateAnnotation)
	})

	router.Delete("/bookmarks/{bookmarkId}", r.Handler.DeleteBookmark)

	router.Route("/reader-config", func(sub chi.Router) {
		sub.Get("/", r.Handler.GetReaderConfig)
		sub.Post("/", r.Handler.SaveReaderConfig)
	})
}
