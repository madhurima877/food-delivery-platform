package repository

import "log"

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (repo *NotificationRepository) SendNotification(userID string, message string) {
	log.Println("Send Notification:", userID, "---", message)
}
