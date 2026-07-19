package handler

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/madhurima877/food-delivery-platform/proto/order"
	orderKafka "github.com/madhurima877/food-delivery-platform/services/order-service/kafka"
	"github.com/madhurima877/food-delivery-platform/services/order-service/repository"
)

type OrderHandler struct {
	repo *repository.OrderRepository
	pb.UnimplementedOrderServiceServer
	producer *orderKafka.Producer
}

func NewOrderHandler(repo *repository.OrderRepository, producer *orderKafka.Producer) *OrderHandler {
	return &OrderHandler{repo: repo, producer: producer}
}
func (h *OrderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	fmt.Println("Delete Order Called")
	isDeleted, err := h.repo.DeleteOrder(req.OrderId)
	if err != nil {
		return nil, err
	}
	if !isDeleted {
		return &pb.DeleteOrderResponse{
			Message: "No Order Deleted",
		}, nil
	}

	return &pb.DeleteOrderResponse{
		Message: "Order Deleted",
	}, nil
}

func (h *OrderHandler) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {

	fmt.Println("CreateOrder called")
	fmt.Println(req.CustomerId)
	fmt.Println(req.RestaurantId)
	fmt.Println(req.ProductId)
	fmt.Println(req.Quantity)

	id, err := h.repo.CreateOrder(req.CustomerId, req.RestaurantId)
	if err != nil {
		return nil, err
	}

	idStr := strconv.FormatInt(id, 10)

	err = h.producer.PublishOrderCreated(
		idStr,
		req.CustomerId,
		req.RestaurantId,
		req.ProductId,
		req.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		OrderId: idStr,
		Status:  "CREATED",
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	fmt.Println("GetOrder Called")
	fmt.Println(req.OrderId)
	Order, err := h.repo.GetOrder(req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{
		OrderId:      Order.Id,
		CustomerId:   Order.CustomerId,
		RestaurantId: Order.RestaurantId,
		Status:       Order.Status,
	}, nil

}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	fmt.Println("UpdateOrder Called")
	fmt.Println(req.OrderId)
	Order, err := h.repo.UpdateOrder(req.OrderId, req.Status)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderStatusResponse{
		OrderId:      Order.Id,
		CustomerId:   Order.CustomerId,
		RestaurantId: Order.RestaurantId,
		Status:       Order.Status,
	}, nil
}
