package data

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"serica-go/internal/conf"
)

var Module = fx.Module("data",
	fx.Provide(NewDB, NewRedis, NewRepos, userRepo, bookRepo, bookKeyRepo, searchHistoryRepo),
)

func NewDB(lc fx.Lifecycle, bc *conf.Bootstrap, zlog *zap.Logger) (*Data, error) {
	d, cleanup, err := NewData(bc, zlog)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cleanup()
			return nil
		},
	})
	return d, nil
}

func NewRedis(lc fx.Lifecycle, bc *conf.Bootstrap, zlog *zap.Logger) (*RedisClient, error) {
	client, err := NewRedisClient(bc)
	if err != nil {
		zlog.Warn("Redis not available, running without cache", zap.Error(err))
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error { return nil },
		})
		return nil, nil
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})
	return client, nil
}

type Repos struct {
	User          *UserRepo
	Book          *BookRepo
	BookKey       *BookKeyRepo
	SearchHistory *SearchHistoryRepo
}

func NewRepos(d *Data) *Repos {
	r := newRepos(d.DB)
	return &Repos{
		User:          r.User,
		Book:          r.Book,
		BookKey:       r.BookKey,
		SearchHistory: r.SearchHistory,
	}
}

func userRepo(r *Repos) *UserRepo                   { return r.User }
func bookRepo(r *Repos) *BookRepo                   { return r.Book }
func bookKeyRepo(r *Repos) *BookKeyRepo             { return r.BookKey }
func searchHistoryRepo(r *Repos) *SearchHistoryRepo { return r.SearchHistory }
