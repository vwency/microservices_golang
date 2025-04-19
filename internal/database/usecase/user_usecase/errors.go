package user_usecase

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound      = status.Error(codes.NotFound, "user not found")
	ErrUserAlreadyExists = status.Error(codes.AlreadyExists, "user already exists")
)
