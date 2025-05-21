package transport

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
)

func GRPCErrorWrapper(err error) error {
	if err == nil {
		return nil
	}

	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	wrappedErr := endpoints.WrapServiceError(err)
	if wrappedErr != nil {
		return status.Error(wrappedErr.Code, wrappedErr.Message)
	}

	return status.Errorf(codes.Internal, "%v", err)
}
