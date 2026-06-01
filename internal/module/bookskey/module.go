package bookskey

import "go.uber.org/fx"

var Module = fx.Module("module-bookskey",
	fx.Provide(NewHandler, NewRoutes),
)
