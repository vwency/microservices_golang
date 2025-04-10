package user_usecase

import (
	"github.com/vwency/microservices_golang/internal/database/repository/user_repository"
	"go.uber.org/zap"
)

type Dependencies struct {
	UserRepo user_repository.UserRepository
	Logger   *zap.Logger
}
