package handler

import (
	"context"
	"fmt"

	pb "github.com/madhurima877/food-delivery-platform/proto/payment"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/repository"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	repo *repository.PaymentRepository
}

func NewPaymentHandler(repo *repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{repo: repo}
}

func (h *PaymentHandler) ProcessPayment(
	ctx context.Context,
	req *pb.ProcessPaymentRequest,
) (*pb.ProcessPaymentResponse, error) {

	fmt.Println("Process Payment Called")

	isCompleted, price, err := h.repo.ProcessPayment(
		req.OrderId,
		req.Price,
		req.UserId,
	)
	if err != nil {
		return nil, err
	}

	if isCompleted {
		return &pb.ProcessPaymentResponse{
			Status: "PaymentDone",
			Price:  price,
		}, nil
	}

	return &pb.ProcessPaymentResponse{
		Status: "PaymentFailed",
		Price:  price,
	}, nil
}
