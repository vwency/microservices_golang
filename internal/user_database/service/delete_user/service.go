package delete_user

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vwency/microservices_golang/internal/user_database/repository/user_repository"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

	type Service interface {
		DeleteUser(ctx context.Context, req Request) (Response, error)
	}

	type service struct {
		repo   user_repository.UserRepository
		logger log.Logger
	}

	func NewService(repo user_repository.UserRepository, logger log.Logger) Service {
		return &service{
			repo:   repo,
			logger: log.With(logger, "service", "delete_user"),
		}
	}

	type Request struct {
		UserID string
	}

	type Response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	func (s *service) DeleteUser(ctx context.Context, req Request) (Response, error) {
		logger := log.With(s.logger, "method", "DeleteUser")

		if req.UserID == "" {
			level.Error(logger).Log("msg", "user_id is required")
			return Response{Success: false}, error_hndl.NewError(codes.InvalidArgument, "user_id is required")
		}

		_, err := s.repo.GetUserByID(req.UserID)
		if err != nil {
			if isInvalidUUIDError(err) {
				level.Warn(logger).Log("msg", "invalid user ID format", "user_id", req.UserID)
				return Response{Success: false}, error_hndl.NewError(codes.InvalidArgument, "invalid user ID format, must be valid UUID")
			}
			if status.Code(err) == codes.NotFound {
				level.Warn(logger).Log("msg", "user not found", "user_id", req.UserID)
				return Response{Success: false}, error_hndl.NewError(codes.NotFound, "user not found")
			}
			level.Error(logger).Log("msg", "failed to get user", "err", err)
			return Response{Success: false}, error_hndl.NewError(codes.Internal, "failed to get user: "+err.Error())
		}

		err = s.repo.DeleteUser(req.UserID)
		if err != nil {
			if isInvalidUUIDError(err) {
				level.Warn(logger).Log("msg", "invalid user ID format on delete", "user_id", req.UserID)
				return Response{Success: false}, error_hndl.NewError(codes.InvalidArgument, "invalid user ID format, must be valid UUID")
			}
			level.Error(logger).Log("msg", "failed to delete user", "err", err)
			return Response{Success: false}, error_hndl.NewError(codes.Internal, "failed to delete user: "+err.Error())
		}

		level.Info(logger).Log("msg", "user deleted successfully", "user_id", req.UserID)
		return Response{
			Success: true,
			Message: "user deleted successfully",
		}, nil
	}

	func isInvalidUUIDError(err error) bool {
		return strings.Contains(err.Error(), "SQLSTATE 22P02") ||
			strings.Contains(err.Error(), "invalid input syntax for type uuid")
	}
	