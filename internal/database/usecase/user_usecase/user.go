package user_usecase

import (
	"errors"
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/models"
)

func (uc *InitUseCase) GetUser(username, email string) (*models.User, error) {
	user, err := uc.repo.GetUserByUsernameOrEmail(username, email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return user, nil
}

func (uc *InitUseCase) AddUser(username, password, hashedRt, hashedAt, email string) error {
	existingUser, err := uc.repo.GetUserByUsernameOrEmail(username, email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	if existingUser != nil {
		return errors.New("user with the same username or email already exists")
	}

	user := &models.User{
		Username: username,
		Password: password,
		HashedRt: hashedRt,
		HashedAt: hashedAt,
		Email:    &email,
	}

	if err := uc.repo.AddUser(user); err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

func (uc *InitUseCase) UpdateUserTokens(username, hashedRt, hashedAt string) error {
	user, err := uc.GetUser(username, "")
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user '%s' not found", username)
	}

	err = uc.repo.UpdateUserTokens(user.Username, hashedRt, hashedAt)
	if err != nil {
		return fmt.Errorf("failed to update user tokens: %w", err)
	}

	return nil
}
