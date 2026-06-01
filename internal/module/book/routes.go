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
	router.Get("/books", r.Handler.List)
	router.Get("/books/{id}", r.Handler.GetByID)
	router.Get("/books/search", r.Handler.Search)

	router.Get("/categories", r.Handler.Categories)
	router.Get("/categories/all", r.Handler.CategoriesAll)
	router.Get("/categories/{id}", r.Handler.CategoriesGetByID)
	router.Get("/categories/{id}/books", r.Handler.CategoryBooks)

	router.Get("/authors", r.Handler.Authors)
	router.Get("/authors/{name}/books", r.Handler.AuthorBooks)

	router.Get("/publishers", r.Handler.Publishers)
	router.Get("/publishers/{name}/books", r.Handler.PublisherBooks)

	router.Get("/selection", r.Handler.Selection)

	router.Post("/trans", r.Handler.Trans)
	router.Post("/text-to-speech", r.Handler.TextToSpeech)

	router.Get("/book-notes", r.Handler.BookNotesList)
	router.Post("/book-notes", r.Handler.BookNotesCreate)
	router.Patch("/book-notes/{id}", r.Handler.BookNotesUpdate)
	router.Delete("/book-notes/{id}", r.Handler.BookNotesDelete)
	router.Post("/book-notes/del/{id}", r.Handler.BookNotesDeleteByPost)

	router.Get("/reading-pos", r.Handler.GetReadingPos)
	router.Post("/reading-pos", r.Handler.ReportReadingPos)
}
