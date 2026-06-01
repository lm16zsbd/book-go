package server

import (
	"go.uber.org/fx"

	"serica-go/internal/module/book"
	"serica-go/internal/module/bookskey"
	"serica-go/internal/module/home"
	"serica-go/internal/module/reader"
	"serica-go/internal/module/search"
	"serica-go/internal/module/user"
)

var Module = fx.Module("server",
	user.Module,
	book.Module,
	bookskey.Module,
	home.Module,
	reader.Module,
	search.Module,
	fx.Provide(
		NewRouter,
		NewHTTPServer,
	),
	fx.Invoke(StartServer),
)
