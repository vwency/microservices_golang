package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type GetUserRequest struct {
	UserID   string
	Username string
	Email    string
}

type GetUserResponse struct {
	Found              bool   `json:"found"`
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	HashedPassword     string `json:"hashed_password"`
	HashedRefreshToken string `json:"hashed_refresh_token"`
	HashedAccessToken  string `json:"hashed_access_token"`
	Error              error  `json:"error,omitempty"`
}

func MakeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := s.GetUser(ctx, service.GetUserRequest{
			UserID:   req.UserID,
			Username: req.Username,
			Email:    req.Email,
		})
		if err != nil {
			return nil, WrapServiceError(err)
		}

		return GetUserResponse{
			Found:              user.Found,
			UserID:             user.UserID,
			Username:           user.Username,
			Email:              user.Email,
			HashedPassword:     user.HashedPassword,
			HashedRefreshToken: user.HashedRefreshToken,
			HashedAccessToken:  user.HashedAccessToken,
		}, nil
	}
}
