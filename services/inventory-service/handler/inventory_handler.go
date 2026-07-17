package handler

import (
	"context"
	"fmt"

	pb "github.com/madhurima877/food-delivery-platform/proto/inventory"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/repository"
)

type InventoryHandler struct {
	pb.UnimplementedInventoryServiceServer
	repo *repository.InventoryRepository
}

func NewInventoryHandler(repo *repository.InventoryRepository) *InventoryHandler {
	return &InventoryHandler{repo: repo}
}

func (h *InventoryHandler) ReserveStock(
	ctx context.Context,
	req *pb.ReserveStockRequest,
) (*pb.ReserveStockResponse, error) {

	fmt.Println("ReserveStock called")
	fmt.Println(req.OrderId)
	fmt.Println(req.ProductId)
	fmt.Println(req.Quantity)

	isUpdated, leftStock, err := h.repo.ReserveStock(
		req.ProductId,
		req.Quantity,
	)

	if err != nil {
		return nil, err
	}

	if !isUpdated {
		return &pb.ReserveStockResponse{
			Status:    "NOT_ENOUGH_STOCK",
			LeftStock: 0,
		}, nil
	}

	return &pb.ReserveStockResponse{
		Status:    "RESERVED",
		LeftStock: leftStock,
	}, nil
}
