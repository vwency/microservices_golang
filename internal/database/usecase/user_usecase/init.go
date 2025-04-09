package user_usecase

import (
	"github.com/vwency/microservices_golang/internal/database/repository/user_repository"
)

type InitUseCase struct {
	repo user_repository.UserRepository
}

func NewInitUseCase(repo user_repository.UserRepository) *InitUseCase {
	return &InitUseCase{repo: repo}
}
