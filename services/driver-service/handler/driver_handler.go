package handler

import (
	"context"

	pb "github.com/madhurima877/food-delivery-platform/proto/driver"
	"github.com/madhurima877/food-delivery-platform/services/driver-service/repository"
)

type DriverHandler struct {
	pb.UnimplementedDriverServiceServer
	repo *repository.DriverRepository
}

func NewDriverHandler(repo *repository.DriverRepository) *DriverHandler {
	return &DriverHandler{
		repo: repo,
	}
}

func (h *DriverHandler) AssignDriver(ctx context.Context, req *pb.DriverRequest) (*pb.DriverResponse, error) {
	h.repo.AssignDriver(req.OrderId, req.DriverId)
	return &pb.DriverResponse{
		Status: "Assigned Successfully",
	}, nil
}
