package endpoints

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type AddUserEndpoint struct{ svc service.Service }

func (e *AddUserEndpoint) Handle(ctx context.Context, req service.AddUserRequest) (service.AddUserResponse, error) {
	return e.svc.AddUser(ctx, req)
}
