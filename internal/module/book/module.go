package book

import "go.uber.org/fx"

var Module = fx.Module("module-book",
	fx.Provide(NewHandler, NewRoutes),
)
