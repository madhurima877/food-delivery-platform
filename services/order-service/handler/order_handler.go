package handler

import (
	"context"
	"fmt"

	pb "github.com/madhurima877/food-delivery-platform/proto/order"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	fmt.Println("CreateOrder called")

	return &pb.CreateOrderResponse{
		OrderId: "1",
		Status:  "CREATED",
	}, nil
}
