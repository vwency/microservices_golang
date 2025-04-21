package user_usecase

import (
	"github.com/vwency/microservices_golang/internal/user_database/repository/user_repository"
	"go.uber.org/zap"
)

type UserUsecase struct {
	repo   user_repository.UserRepository
	logger *zap.Logger
}

func New(repo user_repository.UserRepository, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: logger,
	}
}
