package usecase

import (
	"errors"
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/models"
	"github.com/vwency/microservices_golang/internal/database/repository/user_repository"
)

type InitUseCase struct {
	repo user_repository.UserRepository
}

func NewInitUseCase(repo user_repository.UserRepository) *InitUseCase {
	return &InitUseCase{repo: repo}
}

func (uc *InitUseCase) InitDatabase() error {
	err := uc.repo.RunMigrations()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	return nil
}

func (uc *InitUseCase) GetUser(username, email string) (*models.User, error) {
	user, err := uc.repo.GetUserByUsernameOrEmail(username, email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, nil // Explicitly return nil when user not found
	}
	return user, nil
}

func (uc *InitUseCase) AddUser(username, password, hashedRt, hashedAt string) error {
	existingUser, err := uc.repo.GetUserByUsernameOrEmail(username, "")
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	if existingUser != nil {
		return errors.New("user with the same username already exists")
	}

	// Создание нового пользователя
	user := &models.User{
		Username: username,
		Password: password,
		HashedRt: hashedRt,
		HashedAt: hashedAt,
	}

	// Добавление пользователя в базу данных
	err = uc.repo.AddUser(user)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

func (uc *InitUseCase) UpdateUserTokens(username, hashedRt, hashedAt string) error {
	// Проверка, существует ли пользователь
	user, err := uc.GetUser(username, "")
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Обновление токенов пользователя
	err = uc.repo.UpdateUserTokens(user.Username, hashedRt, hashedAt)
	if err != nil {
		return fmt.Errorf("failed to update user tokens: %w", err)
	}

	return nil
}
