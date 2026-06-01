package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	_ "serica-go/docs"
	"serica-go/internal/wire"
)

// @title           SericaMind Backend Service
// @version         1.0.0
// @description     This is SericaMind backend api service
// @schemes         http https
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	app := fx.New(
		wire.Module,
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
	)

	app.Run()

	if err := app.Err(); err != nil {
		logger.Fatal("application failed", zap.Error(err))
	}
}
