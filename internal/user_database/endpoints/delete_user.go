package endpoints

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type DeleteUserEndpoint struct{ svc service.Service }

func (e *DeleteUserEndpoint) Handle(ctx context.Context, req service.DeleteUserRequest) (service.DeleteUserResponse, error) {
	return e.svc.DeleteUser(ctx, req)
}
