package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/service/errors"
)

func GRPCErrorWrapper(err error) error {
	if err == nil {
		return nil
	}

	// Handle your custom error type first
	if e, ok := err.(*errors.Error); ok {
		return status.Error(e.Code, e.Message)
	}

	// Then handle endpoint's GRPCError
	if e, ok := err.(*endpoints.GRPCError); ok {
		return status.Error(e.Code, e.Message)
	}

	// Then check for standard gRPC status errors
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Default to internal error
	return status.Errorf(codes.Internal, "%v", err)
}
