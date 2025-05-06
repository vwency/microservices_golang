package service

import (
	"context"
)

func (s *userService) UpdateUser(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error) {
	if req.UserID == "" || len(req.Updates) == 0 {
		return UpdateUserResponse{}, ErrInvalidArgument
	}

	// Валидация обновляемых полей
	for field := range req.Updates {
		switch field {
		case "username", "email", "password":
			continue
		default:
			s.logger.Log("error", "invalid update field", "field", field)
			return UpdateUserResponse{}, ErrInvalidArgument
		}
	}

	// Если обновляется пароль - хешируем его
	if password, ok := req.Updates["password"].(string); ok {
		hashed, err := hashPassword(password)
		if err != nil {
			return UpdateUserResponse{}, ErrInternal
		}
		req.Updates["hashed_password"] = hashed
		delete(req.Updates, "password")
	}

	user, err := s.repo.UpdateUser(ctx, req.UserID, req.Updates)
	if err != nil {
		s.logger.Log("error", "user update failed", "userID", req.UserID, "err", err)
		return UpdateUserResponse{}, ErrInternal
	}

	return UpdateUserResponse{User: user}, nil
}
