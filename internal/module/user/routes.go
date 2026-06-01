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
	router.Post("/user/register", r.Handler.Register)
	router.Get("/user/profile", r.Handler.GetProfile)
	router.Patch("/user/profile", r.Handler.UpdateProfile)

	router.Get("/user/favourites", r.Handler.GetFavourites)
	router.Post("/books/{bookId}/favourite", r.Handler.ToggleFavourite)
	router.Post("/books/cancelFavourite", r.Handler.BatchCancelFavourite)

	router.Get("/books/{bookId}/bookmarks", r.Handler.GetBookmarks)
	router.Post("/books/{bookId}/bookmarks", r.Handler.CreateBookmark)
	router.Delete("/bookmarks/{bookmarkId}", r.Handler.DeleteBookmark)

	router.Get("/books/{bookId}/annotations", r.Handler.GetAnnotations)
	router.Post("/books/{bookId}/annotations", r.Handler.CreateAnnotation)

	router.Get("/reader-config", r.Handler.GetReaderConfig)
	router.Post("/reader-config", r.Handler.SaveReaderConfig)
}
