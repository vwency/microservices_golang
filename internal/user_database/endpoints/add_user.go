package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	"github.com/vwency/microservices_golang/internal/user_database/service/add_user"
)

type AddUserRequest struct {
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
	UserID             string
}

type AddUserResponse struct {
	Success bool
	Message string
}

func MakeAddUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddUserRequest)
		res, err := s.AddUser(ctx, add_user.Request{
			Username:           req.Username,
			Email:              req.Email,
			HashedPassword:     req.HashedPassword,
			HashedRefreshToken: req.HashedRefreshToken,
			HashedAccessToken:  req.HashedAccessToken,
			UserID:             req.UserID,
		})

		if err != nil {
			wrappedErr := WrapServiceError(err)
			return nil, wrappedErr
		}
		return AddUserResponse{Success: res.Success, Message: res.Message}, nil
	}
}
