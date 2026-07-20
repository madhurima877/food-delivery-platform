package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	"github.com/madhurima877/food-delivery-platform/services/driver-service/repository"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	repo   *repository.DriverRepository
}

func NewConsumer(repo *repository.DriverRepository) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "payment.completed",
		GroupID: "driver-group"})
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

			driverID := "DRIVER-101"

			c.repo.AssignDriver(
				event.OrderID,
				driverID,
			)
			log.Println(
				"Driver assigned for order:",
				event.OrderID,
				"DriverID:",
				driverID,
			)
		}
	}
}
