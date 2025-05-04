package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
)

// ValidateRequest defines the request structure.
type ValidateRequest struct {
	AccessToken string
}

// ValidateResponse defines the response structure.
type ValidateResponse struct {
	Valid     bool
	UserId    string
	Roles     []string
	ExpiresAt int64
}

// MakeValidateEndpoint creates the endpoint for validating the access token.
func MakeValidateEndpoint(svc service.ValidateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*ValidateRequest)
		result, err := svc.ValidateAccessToken(ctx, req.AccessToken)
		if err != nil {
			// Enhanced error logging and returning a more specific error message
			return nil, err
		}

		return &ValidateResponse{
			Valid:     result.Valid,
			UserId:    result.UserID,
			Roles:     result.Roles,
			ExpiresAt: result.ExpiresAt,
		}, nil
	}
}
