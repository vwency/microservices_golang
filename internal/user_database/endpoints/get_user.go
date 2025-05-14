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
	Found              bool
	UserID             string
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
	Err                error
}

func (r GetUserResponse) error() error { return r.Err }

func MakeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := s.GetUser(ctx, service.GetUserRequest{
			UserID:   req.UserID,
			Username: req.Username,
			Email:    req.Email,
		})
		if err != nil {
			return GetUserResponse{Err: err}, err
		}
		if !user.Found {
			return GetUserResponse{Found: false}, nil
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
