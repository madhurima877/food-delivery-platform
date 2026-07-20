package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/repository"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	repo     *repository.PaymentRepository
	reader   *kafka.Reader
	producer *Producer
}

func NewConsumer(repo *repository.PaymentRepository, producer *Producer) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "inventory.reserved",
		GroupID: "payment-group",
	})
	return &Consumer{repo: repo, reader: reader, producer: producer}
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

			var event models.InventoryReservedEvent

			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Println("Error decoding inventory.reserved event:", err)
				continue
			}

			log.Println("Inventory reserved event received for order:", event.OrderID)

			isCompleted, price, err := c.repo.ProcessPayment(
				event.OrderID,
				float32(event.TotalPrice),
				event.CustomerID,
			)
			if err != nil {
				log.Println("Payment error:", err)
				continue
			}

			if !isCompleted {
				log.Println("Payment failed for order:", event.OrderID)
				continue
			}
			err = c.producer.PublishEvent(
				event.OrderID,
				event.CustomerID,
				"COMPLETED",
				price,
			)
			if err != nil {
				log.Println("Error publishing payment.completed event:", err)
				continue
			}

			log.Println(
				"Payment completed for order:",
				event.OrderID,
				"price:",
				price,
			)
		}
	}
}
