package transport

import (
	"context"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *grpcServer) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	_, resp, err := s.logout.ServeGRPC(ctx, req)
	if err != nil {
		return nil, GRPCErrorWrapper(err) // Здесь применяем обёртку
	}
	return resp.(*authv1.LogoutResponse), nil
}

func decodeLogoutRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authv1.LogoutRequest)
	if !ok {
		return nil, GRPCErrorWrapper(status.Error(codes.InvalidArgument, "invalid request type"))
	}
	return req, nil
}

func encodeLogoutResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*authv1.LogoutResponse)
	if !ok {
		return nil, GRPCErrorWrapper(status.Error(codes.Internal, "invalid response type"))
	}
	return resp, nil
}
