package endpoints

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type InitDatabaseEndpoint struct{ svc service.Service }

func (e *InitDatabaseEndpoint) Handle(ctx context.Context, req service.InitDatabaseRequest) (service.InitDatabaseResponse, error) {
	return e.svc.InitDatabase(ctx, req)
}
