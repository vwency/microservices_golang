package endpoints

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type GetUserEndpoint struct{ svc service.Service }

func (e *GetUserEndpoint) Handle(ctx context.Context, req service.GetUserRequest) (service.GetUserResponse, error) {
	return e.svc.GetUser(ctx, req)
}
