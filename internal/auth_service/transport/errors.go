package transport

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
)

// GRPCErrorWrapper converts service-layer errors into proper gRPC errors with status codes.
func GRPCErrorWrapper(err error) error {
	if err == nil {
		return nil
	}

	// Если это уже gRPC статус - возвращаем как есть
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Используем вашу обёртку из endpoints для преобразования
	wrappedErr := endpoints.WrapServiceError(err)

	// Если после обёртки получаем *endpoints.GRPCError, создаём status.Error
	var grpcErr *endpoints.GRPCError
	if errors.As(wrappedErr, &grpcErr) {
		return status.Error(grpcErr.Code, grpcErr.Message)
	}

	// В крайнем случае возвращаем Internal с сообщением
	return status.Errorf(codes.Internal, "unknown error: %v", err)
}

// Helper для форматирования логов (опционально)
func LogGRPCError(err error) string {
	if st, ok := status.FromError(err); ok {
		return fmt.Sprintf("gRPC status error - Code: %s, Message: %s", st.Code(), st.Message())
	}
	return fmt.Sprintf("non gRPC error: %v", err)
}
