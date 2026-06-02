package book

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
	router.Route("/books", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.List))
		sub.Get("/{id}", httputil.Wrap(r.Handler.GetByID))
	})

	router.Route("/categories", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.Categories))
		sub.Get("/all", httputil.Wrap(r.Handler.CategoriesAll))
		sub.Get("/{id}", httputil.Wrap(r.Handler.CategoriesGetByID))
		sub.Get("/{id}/books", httputil.Wrap(r.Handler.CategoryBooks))
	})

	router.Route("/authors", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.Authors))
		sub.Get("/{name}/books", httputil.Wrap(r.Handler.AuthorBooks))
	})

	router.Route("/publishers", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.Publishers))
		sub.Get("/{name}/books", httputil.Wrap(r.Handler.PublisherBooks))
	})

	router.Route("/book-notes", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.BookNotesList))
		sub.Post("/", httputil.Wrap(r.Handler.BookNotesCreate))
		sub.Patch("/{id}", httputil.Wrap(r.Handler.BookNotesUpdate))
		sub.Delete("/{id}", httputil.Wrap(r.Handler.BookNotesDelete))
		sub.Post("/del/{id}", httputil.Wrap(r.Handler.BookNotesDeleteByPost))
	})

	router.Route("/reading-pos", func(sub chi.Router) {
		sub.Get("/", httputil.Wrap(r.Handler.GetReadingPos))
		sub.Post("/", httputil.Wrap(r.Handler.ReportReadingPos))
	})

	router.Get("/search", httputil.Wrap(r.Handler.Search))
	router.Get("/selection", httputil.Wrap(r.Handler.Selection))
	router.Post("/trans", httputil.Wrap(r.Handler.Trans))
	router.Post("/text-to-speech", httputil.Wrap(r.Handler.TextToSpeech))
}
