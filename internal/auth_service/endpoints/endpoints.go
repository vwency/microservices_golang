package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	"go.opentelemetry.io/otel"
)

type Endpoints struct {
	Login               endpoint.Endpoint
	Logout              endpoint.Endpoint
	Register            endpoint.Endpoint
	Refresh             endpoint.Endpoint
	ValidateAccessToken endpoint.Endpoint
}

func MakeEndpoints(s service.AuthService) Endpoints {
	tracer := otel.Tracer("auth_service")

	return Endpoints{
		Login:               TraceEndpoint(tracer, "Login")(MakeLoginEndpoint(s)),
		Logout:              TraceEndpoint(tracer, "Logout")(MakeLogoutEndpoint(s)),
		Register:            TraceEndpoint(tracer, "Register")(MakeRegisterEndpoint(s)),
		Refresh:             TraceEndpoint(tracer, "Refresh")(MakeRefreshEndpoint(s)),
		ValidateAccessToken: TraceEndpoint(tracer, "ValidateAccessToken")(MakeValidateEndpoint(s)),
	}
}
