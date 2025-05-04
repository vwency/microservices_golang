package transport

import (
	"context"

	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func decodeRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*authv1.RegisterRequest)
	return &service.RegisterRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}, nil
}

func encodeRegisterResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*service.RegisterResponse)
	return &authv1.RegisterResponse{
		UserId:       resp.UserID,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
	}, nil
}
