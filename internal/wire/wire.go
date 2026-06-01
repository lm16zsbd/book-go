package wire

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"serica-go/internal/conf"
	"serica-go/internal/data"
	"serica-go/internal/server"
)

var Module = fx.Options(
	fx.Provide(NewLogger),
	conf.Module,
	data.Module,
	server.Module,
)

func NewLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}
