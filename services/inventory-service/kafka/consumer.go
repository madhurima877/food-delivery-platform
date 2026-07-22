package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/repository"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader        *kafka.Reader
	repo          *repository.InventoryRepository
	producer      *Producer
	retryReader   *kafka.Reader
	restoreReader *kafka.Reader
}

func NewConsumer(repo *repository.InventoryRepository, writer *Producer) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "order.created",
		GroupID: "inventory-group",
	})
	retryReader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{"localhost:19092"},
			Topic:   "order.created.retry",
			GroupID: "inventory-retry-group",
		})
	restoreReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "payment.failed",
		GroupID: "inventory-payment-failed-group",
	})
	return &Consumer{reader: reader, repo: repo, producer: writer, retryReader: retryReader, restoreReader: restoreReader}
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

			retryCount := getRetryCount(msg)
			var event models.OrderCreatedEvent
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Println("Error decoding event:", err)
				continue
			}

			status, leftStock, err := c.repo.ReserveStock(
				event.OrderID,
				event.ProductID,
				event.Quantity,
			)
			if err != nil {
				log.Println("Error reserving stock:", err)
				retryCount++

				err = c.producer.PublishRetryEvent(msg.Value, retryCount)
				if err != nil {
					log.Println("Error publishing retry event:", err)
				}
				continue
			}

			if status == "DUPLICATE" {
				log.Println("Duplicate event ignored for order:", event.OrderID)
				continue
			}

			if status == "NOT_ENOUGH_STOCK" {
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

func getRetryCount(msg kafka.Message) int {
	for _, header := range msg.Headers {
		if header.Key == "retry-count" {
			cntstr := string(header.Value)
			cnt, _ := strconv.Atoi(cntstr)
			return cnt
		}

	}
	return 0

}

func (c *Consumer) ReadRetryConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.retryReader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}
			retryCount := getRetryCount(msg)
			log.Println("Processing retry event, attempt:", retryCount)
			var event models.OrderCreatedEvent

			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Println("Error decoding retry event:", err)
				continue
			}
			status, leftStock, err := c.repo.ReserveStock(
				event.OrderID,
				event.ProductID,
				event.Quantity,
			)
			if err != nil {
				log.Println("Error reserving stock:", err)
				if retryCount >= 3 {
					log.Println("Max retries reached for order:", event.OrderID)

					err = c.producer.PublishDLQEvent(msg.Value)
					if err != nil {
						log.Println("Error publishing event to DLQ:", err)
						continue
					}

					log.Println("Event sent to DLQ for order:", event.OrderID)
					continue
				}

				retryCount++

				err = c.producer.PublishRetryEvent(msg.Value, retryCount)
				if err != nil {
					log.Println("Error publishing retry event:", err)
				}

				continue
			}

			if status == "DUPLICATE" {
				log.Println("Duplicate event ignored for order:", event.OrderID)
				continue
			}

			if status == "NOT_ENOUGH_STOCK" {
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

func (c *Consumer) ReadRestoreConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.restoreReader.ReadMessage(ctx)
			if err != nil {
				continue
			}
			var event models.PaymentFailedEvent
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				continue
			}
			log.Println("Payment failed event received")

			err = c.repo.RestoreStock(event.ProductID, event.Quantity)
			if err != nil {
				continue
			}
		}
	}
}
