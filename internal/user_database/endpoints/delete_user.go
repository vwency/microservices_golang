package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type DeleteUserRequest struct {
	UserID string
}

type DeleteUserResponse struct {
	Success bool
	Message string
}

func MakeDeleteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		res, err := s.DeleteUser(ctx, service.DeleteUserRequest{
			UserID: req.UserID,
		})
		if err != nil {
			return nil, WrapServiceError(err)
		}
		return DeleteUserResponse{
			Success: res.Success,
			Message: "User deleted successfully",
		}, nil
	}
}
