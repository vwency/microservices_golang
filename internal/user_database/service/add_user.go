package service

import (
	"context"
)

func (s *userService) AddUser(ctx context.Context, req AddUserRequest) (AddUserResponse, error) {
	// Хеширование пароля
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		s.logger.Log("error", "password hashing failed", "err", err)
		return AddUserResponse{}, ErrInternal
	}

	user := User{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	// Проверка уникальности пользователя
	if exists, _ := s.repo.UserExists(ctx, req.Username, req.Email); exists {
		s.logger.Log("error", "user already exists")
		return AddUserResponse{}, ErrAlreadyExists
	}

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Log("error", "user creation failed", "err", err)
		return AddUserResponse{}, ErrInternal
	}

	s.logger.Log("msg", "user created", "userID", userID)
	return AddUserResponse{UserID: userID}, nil
}

func hashPassword(password string) (string, error) {
	// Реальная реализация должна использовать bcrypt или аналоги
	return "hashed_" + password, nil
}
