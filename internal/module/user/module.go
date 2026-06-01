package user

import "go.uber.org/fx"

var Module = fx.Module("module-user",
	fx.Provide(NewHandler, NewRoutes),
)
