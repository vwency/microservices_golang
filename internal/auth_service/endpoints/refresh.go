package endpoints

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	"github.com/vwency/microservices_golang/proto/auth_service"
)

// MakeRefreshEndpoint creates an endpoint for refresh
func MakeRefreshEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*auth_service.RefreshRequest)
		resp, err := svc.Refresh(ctx, req.RefreshToken) // Pass the refreshToken from the request
		if err != nil {
			return nil, fmt.Errorf("refresh failed: %w", err)
		}

		return resp, nil
	}
}
