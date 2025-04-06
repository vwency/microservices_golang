package usecase

import (
	"github.com/vwency/microservices_golang/internal/database_init_service/repository"
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

	err := uc.repo.RunMigrations()
	if err != nil {
		return err
	}

	logger.Info("Database initialized successfully")
	return nil
}
