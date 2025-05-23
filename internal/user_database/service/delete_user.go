package service

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteUserRequest struct {
	UserID string
}

type DeleteUserResponse struct {
	Success bool
	Message string
}

func (s *userService) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	if req.UserID == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "user_id is required")
		level.Error(logger).Log("msg", err.Error())
		return DeleteUserResponse{Success: false, Message: err.Message}, err
	}

	_, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		if isInvalidUUIDError(err) {
			errInvalid := error_hndl.NewError(codes.InvalidArgument, "invalid user ID format, must be valid UUID")
			level.Warn(logger).Log("msg", errInvalid.Error(), "userID", req.UserID)
			return DeleteUserResponse{Success: false, Message: errInvalid.Message}, errInvalid
		}
		if status.Code(err) == codes.NotFound {
			errNotFound := error_hndl.NewError(codes.NotFound, "user not found")
			level.Warn(logger).Log("msg", errNotFound.Error(), "userID", req.UserID)
			return DeleteUserResponse{Success: false, Message: errNotFound.Message}, errNotFound
		}
		errInternal := error_hndl.NewError(codes.Internal, "failed to get user: "+err.Error())
		level.Error(logger).Log("msg", errInternal.Error(), "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false, Message: errInternal.Message}, errInternal
	}

	err = s.repo.UserRepo.DeleteUser(req.UserID)
	if err != nil {
		if isInvalidUUIDError(err) {
			errInvalid := error_hndl.NewError(codes.InvalidArgument, "invalid user ID format, must be valid UUID")
			level.Warn(logger).Log("msg", errInvalid.Error(), "userID", req.UserID)
			return DeleteUserResponse{Success: false, Message: errInvalid.Message}, errInvalid
		}
		errInternal := error_hndl.NewError(codes.Internal, "failed to delete user: "+err.Error())
		level.Error(logger).Log("msg", errInternal.Error(), "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false, Message: errInternal.Message}, errInternal
	}

	level.Info(logger).Log("msg", "user deleted successfully", "userID", req.UserID)
	return DeleteUserResponse{
		Success: true,
		Message: "user deleted successfully",
	}, nil
}

func isInvalidUUIDError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 22P02") ||
		strings.Contains(err.Error(), "invalid input syntax for type uuid")
}
