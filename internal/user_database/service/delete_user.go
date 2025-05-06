package service

import (
	"context"
)

func (s *userService) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	if req.UserID == "" {
		return DeleteUserResponse{Success: false}, ErrInvalidArgument
	}

	// Проверяем существование пользователя
	if _, err := s.repo.GetUserByID(ctx, req.UserID); err != nil {
		s.logger.Log("error", "user not found for deletion", "userID", req.UserID)
		return DeleteUserResponse{Success: false}, ErrNotFound
	}

	err := s.repo.UserRepo.(ctx, req.UserID)
	if err != nil {
		s.logger.Log("error", "user deletion failed", "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false}, ErrInternal
	}

	s.logger.Log("msg", "user deleted", "userID", req.UserID)
	return DeleteUserResponse{Success: true}, nil
}
