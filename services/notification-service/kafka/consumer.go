package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	"github.com/madhurima877/food-delivery-platform/services/notification-service/repository"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	repo   *repository.NotificationRepository
}

func NewConsumer(repo *repository.NotificationRepository) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "payment.completed",
		GroupID: "notification-group"})
	return &Consumer{
		reader: reader,
		repo:   repo,
	}
}

func (c *Consumer) ReadConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}

				log.Println("Kafka Read Error:", err)
				continue
			}

			var event models.PaymentCompletedEvent

			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Println("Error decoding payment.completed event:", err)
				continue
			}

			c.repo.SendNotification(
				event.CustomerID,
				"Payment completed successfully for order "+event.OrderID,
			)

			log.Println("Notification sent for order:", event.OrderID)
		}
	}
}
