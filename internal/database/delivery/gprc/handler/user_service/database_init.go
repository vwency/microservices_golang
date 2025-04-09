package handler_user_service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
)

type DatabaseInitHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	usecase *usecase.InitUseCase
}

func NewDatabaseInitHandler(uc *usecase.InitUseCase) *DatabaseInitHandler {
	return &DatabaseInitHandler{
		usecase: uc,
	}
}

// InitDatabase выполняет инициализацию базы данных
func (h *DatabaseInitHandler) InitDatabase(ctx context.Context, req *pb.InitRequest) (*pb.InitResponse, error) {
	err := h.usecase.InitDatabase()
	if err != nil {
		return &pb.InitResponse{
			Success: false,
			Message: "Initialization failed: " + err.Error(),
		}, nil
	}

	return &pb.InitResponse{
		Success: true,
		Message: "Database initialized successfully",
	}, nil
}
