package handler

import (
	"context"

	pb "github.com/vwency/microservices_golang/proto/database"

	"github.com/vwency/microservices_golang/internal/database/usecase"
)

type DatabaseInitHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	usecase usecase.DatabaseInit
}

func NewHandler(uc usecase.DatabaseInit) *DatabaseInitHandler {
	return &DatabaseInitHandler{
		usecase: uc,
	}
}

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
