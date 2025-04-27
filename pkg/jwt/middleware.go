package jwt

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "Bearer "
)

// ContextKey тип для ключей контекста
type ContextKey string

const (
	// ClaimsKey ключ для хранения claims в контексте
	ClaimsKey ContextKey = "claims"
)

// JWTAuthInterceptor создает gRPC interceptor для JWT аутентификации
func JWTAuthInterceptor(jwtManager *JWTManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Пропускаем аутентификацию для методов AuthService
		if strings.Contains(info.FullMethod, "AuthService") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		authHeaders := md.Get(authorizationHeader)
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		token := strings.TrimPrefix(authHeaders[0], bearerPrefix)
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		newCtx := context.WithValue(ctx, ClaimsKey, claims)
		return handler(newCtx, req)
	}
}

func GetClaimsFromContext(ctx context.Context) (map[string]interface{}, bool) {
	claims, ok := ctx.Value(ClaimsKey).(map[string]interface{})
	return claims, ok
}
