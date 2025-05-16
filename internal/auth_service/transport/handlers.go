package transport

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
)

func makeLoginHandler(ep endpoints.Endpoints) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep.Login,
		decodeLoginRequest,
		encodeLoginResponse,
	)
}

func makeLogoutHandler(ep endpoints.Endpoints) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep.Logout,
		decodeLogoutRequest,
		encodeLogoutResponse,
	)
}

func makeRegisterHandler(ep endpoints.Endpoints) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep.Register,
		decodeRegisterRequest,
		encodeRegisterResponse,
	)
}

func makeRefreshHandler(ep endpoints.Endpoints) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep.Refresh,
		decodeRefreshRequest,
		encodeRefreshResponse,
	)
}

func makeValidateHandler(ep endpoints.Endpoints) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep.ValidateAccessToken,
		decodeValidateRequest,
		encodeValidateResponse,
	)
}
