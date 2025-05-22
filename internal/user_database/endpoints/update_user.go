package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type UpdateUserRequest struct {
	UserID             string
	HashedRefreshToken string
	HashedAccessToken  string
}

type UpdateUserResponse struct {
	Success bool
	Message string
}

func MakeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserRequest)
		res, err := s.UpdateUser(ctx, service.UpdateUserRequest{
			UserID:             req.UserID,
			HashedRefreshToken: req.HashedRefreshToken,
			HashedAccessToken:  req.HashedAccessToken,
		})

		if err != nil {
			return nil, WrapServiceError(err)
		}

		return UpdateUserResponse{
			Success: res.Success,
			Message: res.Message,
		}, nil
	}
}
