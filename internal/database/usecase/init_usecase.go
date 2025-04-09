package usecase

import (
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/repository"
	"github.com/vwency/microservices_golang/pkg/logger"
)

type DatabaseInit interface {
	InitDatabase() error
}

type databaseInitUsecase struct {
	repo repository.DatabaseInitRepository
}

func NewUsecase(repo repository.DatabaseInitRepository) DatabaseInit {
	return &databaseInitUsecase{
		repo: repo,
	}
}

func (uc *databaseInitUsecase) InitDatabase() error {
	logger.Info("Starting database initialization")

	if uc == nil {
		return fmt.Errorf("usecase is nil")
	}
	if uc.repo == nil {
		return fmt.Errorf("repository is not initialized")
	}

	logger.Info("Attempting to run migrations")
	err := uc.repo.RunMigrations()
	if err != nil {
		logger.Error("Migration failed: %v", err)
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("Database initialized successfully")
	return nil
}
