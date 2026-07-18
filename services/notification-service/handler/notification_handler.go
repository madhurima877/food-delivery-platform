package handler

import (
	"context"

	pb "github.com/madhurima877/food-delivery-platform/proto/notification"
	"github.com/madhurima877/food-delivery-platform/services/notification-service/repository"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	h.repo.SendNotification(req.UserId, req.Message)
	return &pb.NotificationResponse{
		Status: "success",
	}, nil
}
