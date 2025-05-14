package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCStatusCoder interface {
	GRPCStatus() *status.Status
}

func ErrorEncoder(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			resp, err := next(ctx, request)
			if err == nil {
				return resp, nil
			}

			logger.Log("error", err)

			// If the error already implements GRPCStatusCoder, use it
			if sc, ok := err.(GRPCStatusCoder); ok {
				return nil, sc.GRPCStatus().Err()
			}

			// If it's a gRPC status error, return it directly
			if _, ok := status.FromError(err); ok {
				return nil, err
			}

			// Default to Internal error for unknown errors
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
}
