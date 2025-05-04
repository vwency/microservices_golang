package transport

import (
	"context"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func decodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authv1.LoginRequest)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid request type")
	}
	return req, nil
}

func encodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*authv1.LoginResponse)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid response type")
	}
	return resp, nil
}
