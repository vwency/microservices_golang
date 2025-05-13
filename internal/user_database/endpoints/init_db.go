package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type InitDatabaseRequest struct {
	ConfigPath string
}

type InitDatabaseResponse struct {
	Success bool
	Message string
}

func MakeInitDatabaseEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(InitDatabaseRequest)
		res, err := s.InitDatabase(ctx, service.InitDatabaseRequest{
			ConfigPath: req.ConfigPath,
		})
		if err != nil {
			return InitDatabaseResponse{
				Success: false,
				Message: err.Error(),
			}, err
		}
		return InitDatabaseResponse{
			Success: res.Success,
			Message: "Database initialized successfully",
		}, nil
	}
}
