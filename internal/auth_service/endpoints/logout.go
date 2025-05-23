package endpoints

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeLogoutEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("[DEBUG] Entering MakeLogoutEndpoint handler")

		req, ok := request.(*authv1.LogoutRequest)
		if !ok {
			fmt.Printf("[ERROR] Invalid request type: %T\n", request)
			return nil, status.Error(codes.InvalidArgument, "invalid request type")
		}

		fmt.Printf("[DEBUG] Logout request for username: %s\n", req.Username)

		resp, err := s.Logout(ctx, req)
		if err != nil {
			fmt.Printf("[ERROR] Logout failed: %v\n", err)

			if st, ok := status.FromError(err); ok {
				fmt.Printf("[DEBUG] gRPC status error - Code: %s, Message: %s\n",
					st.Code(), st.Message())
				return nil, st.Err()
			} else {
				fmt.Println("[DEBUG] Non-gRPC error type")
				return nil, status.Error(codes.Internal, err.Error())
			}
		}

		fmt.Printf("[DEBUG] Logout successful for username: %s\n", req.Username)
		return resp, nil
	}
}
