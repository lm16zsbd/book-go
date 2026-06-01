package home

import "go.uber.org/fx"

var Module = fx.Module("home",
	fx.Provide(NewHandler, NewRoutes),
)
