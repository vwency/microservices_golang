package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type AddUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// другие поля
}

type AddUserResponse struct {
	User *service.User `json:"user"`
	Err  error         `json:"error,omitempty"`
}

func (r AddUserResponse) error() error { return r.Err }

func MakeAddUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddUserRequest)
		user, err := s.AddUser(ctx, service.AddUserRequest{
			Username: req.Username,
			Email:    req.Email,
			// другие поля
		})
		if err != nil {
			return AddUserResponse{Err: err}, nil
		}
		return AddUserResponse{User: user}, nil
	}
}
