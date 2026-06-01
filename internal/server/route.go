package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"serica-go/docs"
	"serica-go/internal/conf"
	"serica-go/internal/data"
	"serica-go/internal/module/book"
	"serica-go/internal/module/bookskey"
	"serica-go/internal/module/home"
	"serica-go/internal/module/reader"
	"serica-go/internal/module/search"
	"serica-go/internal/module/user"
	svrMiddleware "serica-go/internal/server/middleware"
)

func NewRouter(
	cfg *conf.Bootstrap,
	d *data.Data,
	redis *data.RedisClient,
	userRepo *data.UserRepo,
	userRoutes *user.Routes,
	bookRoutes *book.Routes,
	booksKeyRoutes *bookskey.Routes,
	homeRoutes *home.Routes,
	readerRoutes *reader.Routes,
	searchRoutes *search.Routes,
) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(svrMiddleware.CORS(cfg.Server.CORS))
	r.Use(svrMiddleware.Recovery)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","dbHealthy":%v}`, d.Healthy)
	})

	r.Get("/api/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
	})

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/", http.StatusMovedPermanently)
	})

	r.Get("/api/*", httpSwagger.Handler(
		httpSwagger.URL("/api/doc.json"),
	))

	r.Route("/v1/client", func(r chi.Router) {
		r.Use(svrMiddleware.RequireDB(d.Healthy))
		r.Use(svrMiddleware.Auth(userRepo, redis, svrMiddleware.SchemaClient, svrMiddleware.SchemaPublic))
		userRoutes.Register(r)
		bookRoutes.Register(r)
		booksKeyRoutes.Register(r)
		homeRoutes.Register(r)
		searchRoutes.Register(r)
	})

	r.Route("/v1/client/reader", func(r chi.Router) {
		r.Use(svrMiddleware.RequireDB(d.Healthy))
		r.Use(svrMiddleware.Auth(userRepo, redis, svrMiddleware.SchemaClient, svrMiddleware.SchemaPublic))
		readerRoutes.Register(r)
	})

	_ = redis

	return r
}
