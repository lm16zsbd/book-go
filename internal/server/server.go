package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"serica-go/internal/conf"
)

func NewHTTPServer(bc *conf.Bootstrap, router chi.Router) *http.Server {
	return &http.Server{
		Addr:    bc.Server.HTTP.Addr,
		Handler: router,
	}
}

func StartServer(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Printf("HTTP server listening on %s\n", srv.Addr)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
