package transport

import (
	"context"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func (s *grpcServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	_, resp, err := s.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*authv1.RegisterResponse), nil
}

func decodeRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*authv1.RegisterRequest)
	return req, nil
}

func encodeRegisterResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
