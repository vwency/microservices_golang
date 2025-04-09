package user_usecase

import (
	"github.com/vwency/microservices_golang/internal/database/repository/user_repository"
	"go.uber.org/zap"
)

// Dependencies содержит все зависимости для usecase
type Dependencies struct {
	UserRepo user_repository.UserRepository // Используем интерфейс из user_repository
	Logger   *zap.Logger
	// Другие зависимости при необходимости
}

// New создает новый экземпляр UserUsecase с внедренными зависимостями
