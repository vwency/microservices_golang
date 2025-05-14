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

			// Если ошибка уже имеет gRPC-статус, возвращаем как есть
			if st, ok := status.FromError(err); ok {
				return nil, st.Err()
			}

			// Если ошибка реализует GRPCStatusCoder, используем её
			if sc, ok := err.(GRPCStatusCoder); ok {
				return nil, sc.GRPCStatus().Err()
			}

			// Все остальные ошибки — Internal
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
}
