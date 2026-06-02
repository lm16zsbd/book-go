package book

import "github.com/go-chi/chi/v5"

const BasePath = "/v1/client"

type Routes struct {
	Handler *Handler
}

func NewRoutes(h *Handler) *Routes {
	return &Routes{Handler: h}
}

func (r *Routes) Register(router chi.Router) {
	router.Route("/books", func(sub chi.Router) {
		sub.Get("/", r.Handler.List)
		sub.Get("/{id}", r.Handler.GetByID)
	})

	router.Route("/categories", func(sub chi.Router) {
		sub.Get("/", r.Handler.Categories)
		sub.Get("/all", r.Handler.CategoriesAll)
		sub.Get("/{id}", r.Handler.CategoriesGetByID)
		sub.Get("/{id}/books", r.Handler.CategoryBooks)
	})

	router.Route("/authors", func(sub chi.Router) {
		sub.Get("/", r.Handler.Authors)
		sub.Get("/{name}/books", r.Handler.AuthorBooks)
	})

	router.Route("/publishers", func(sub chi.Router) {
		sub.Get("/", r.Handler.Publishers)
		sub.Get("/{name}/books", r.Handler.PublisherBooks)
	})

	router.Route("/book-notes", func(sub chi.Router) {
		sub.Get("/", r.Handler.BookNotesList)
		sub.Post("/", r.Handler.BookNotesCreate)
		sub.Patch("/{id}", r.Handler.BookNotesUpdate)
		sub.Delete("/{id}", r.Handler.BookNotesDelete)
		sub.Post("/del/{id}", r.Handler.BookNotesDeleteByPost)
	})

	router.Route("/reading-pos", func(sub chi.Router) {
		sub.Get("/", r.Handler.GetReadingPos)
		sub.Post("/", r.Handler.ReportReadingPos)
	})

	router.Get("/search", r.Handler.Search)
	router.Get("/selection", r.Handler.Selection)
	router.Post("/trans", r.Handler.Trans)
	router.Post("/text-to-speech", r.Handler.TextToSpeech)
}
