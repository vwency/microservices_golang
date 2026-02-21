package endpoints

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type UpdateUserEndpoint struct{ svc service.Service }

func (e *UpdateUserEndpoint) Handle(ctx context.Context, req service.UpdateUserRequest) (service.UpdateUserResponse, error) {
	return e.svc.UpdateUser(ctx, req)
}
