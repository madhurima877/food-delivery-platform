package repository

import (
	"database/sql"
	"log"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}
func (repo *NotificationRepository) SendNotification(userid string, message string) {
	log.Println("Send Notification", userid, "---", message)
}
