package endpoints

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type GetUserByIDEndpoint struct{ svc service.Service }

func (e *GetUserByIDEndpoint) Handle(ctx context.Context, req service.GetUserByIDRequest) (service.GetUserByIDResponse, error) {
	return e.svc.GetUserByID(ctx, req)
}
