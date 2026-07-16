package handler

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/madhurima877/food-delivery-platform/proto/order"
	"github.com/madhurima877/food-delivery-platform/services/order-service/repository"
)

type OrderHandler struct {
	repo *repository.OrderRepository
	pb.UnimplementedOrderServiceServer
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	fmt.Println("CreateOrder called")
	fmt.Println(req.CustomerId)
	fmt.Println(req.RestaurantId)
	id, err := h.repo.CreateOrder(req.CustomerId, req.RestaurantId)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		OrderId: strconv.FormatInt(id, 10),
		Status:  "CREATED",
	}, nil
}
