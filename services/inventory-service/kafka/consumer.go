package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/repository"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	repo     *repository.InventoryRepository
	producer *Producer
}

func NewConsumer(repo *repository.InventoryRepository, writer *Producer) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "order.created",
		GroupID: "inventory-group",
	})
	return &Consumer{reader: reader, repo: repo, producer: writer}
}

func (c *Consumer) ReadConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("Inventory consumer stopped")
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
			var event models.OrderCreatedEvent
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Println("Error decoding event:", err)
				continue
			}

			isUpdated, leftStock, err := c.repo.ReserveStock(event.ProductID, event.Quantity)
			if err != nil {
				log.Println("Error reserving stock:", err)
				continue
			}

			if !isUpdated {
				log.Println("Not enough stock for order:", event.OrderID)
				continue
			}
			price, err := c.repo.GetProductPrice(event.ProductID)
			if err != nil {
				log.Println("Error getting price for product:", event.ProductID, err)
				continue
			}

			totalPrice := price * int(event.Quantity)

			err = c.producer.PublishEvent(
				event.OrderID,
				event.CustomerID,
				event.ProductID,
				int32(totalPrice),
				event.Quantity,
			)
			if err != nil {
				log.Println("Error publishing inventory.reserved event:", err)
				continue
			}

			log.Println(
				"Stock reserved for order:",
				event.OrderID,
				"left stock:",
				leftStock,
			)

		}
	}
}
