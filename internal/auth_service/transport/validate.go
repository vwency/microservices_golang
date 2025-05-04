package transport

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServerValidate struct {
	validate grpc.Handler
	authv1.UnimplementedAuthServiceServer
}

func NewGRPCServerValidate(validateEndpoint kitendpoint.Endpoint) authv1.AuthServiceServer {
	return &grpcServerValidate{
		validate: grpc.NewServer(
			validateEndpoint,
			decodeValidateRequest,
			encodeValidateResponse,
		),
	}
}

func (s *grpcServerValidate) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	return resp.(*authv1.ValidateResponse), nil
}

func decodeValidateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authv1.ValidateRequest)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid request type")
	}
	return req, nil
}

func encodeValidateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*authv1.ValidateResponse)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid response type")
	}
	return resp, nil
}
