package reader

import "go.uber.org/fx"

var Module = fx.Module("reader",
	fx.Provide(NewHandler, NewRoutes),
)
