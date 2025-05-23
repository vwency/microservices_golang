package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/vwency/microservices_golang/internal/user_database/repository/user_repository"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/database"
)

func newDatabaseConnection(lc fx.Lifecycle, cfg config.ServiceConfig, logger *zap.Logger) (*gorm.DB, error) {
	db, err := database.NewGORM(cfg.UserDatabase.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user_database: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := user_repository.RunUserMigrations(db); err != nil {
				return fmt.Errorf("failed to run migrations: %w", err)
			}
			logger.Info("Database migrations completed successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})

	return db, nil
}
