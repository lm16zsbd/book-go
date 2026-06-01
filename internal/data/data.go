package data

import (
	"serica-go/internal/conf"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Data struct {
	DB      *gorm.DB
	Healthy bool
}

func NewData(c *conf.Bootstrap, zlog *zap.Logger) (*Data, func(), error) {
	cleanup := func() {}

	db, err := gorm.Open(postgres.Open(c.Data.Database.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		zlog.Warn("PostgreSQL not available, running without database", zap.Error(err))
		return &Data{DB: nil, Healthy: false}, cleanup, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		zlog.Warn("PostgreSQL ping failed, running without database", zap.Error(err))
		return &Data{DB: nil, Healthy: false}, cleanup, nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)

	cleanup = func() {
		sqlDB.Close()
	}

	return &Data{DB: db, Healthy: true}, cleanup, nil
}
