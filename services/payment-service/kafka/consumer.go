package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/lock"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/repository"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	repo     *repository.PaymentRepository
	reader   *kafka.Reader
	producer *Producer
	lock     *lock.RedisLock
}

func NewConsumer(repo *repository.PaymentRepository, producer *Producer, lock *lock.RedisLock) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "inventory.reserved",
		GroupID: "payment-group",
	})
	return &Consumer{repo: repo, reader: reader, producer: producer, lock: lock}
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

			c.processPaymentEvent(ctx, event)

		}
	}
}

func (c *Consumer) processPaymentEvent(ctx context.Context, event models.InventoryReservedEvent) {
	acquired, err := c.lock.Acquire(ctx, event.OrderID)
	if err != nil {
		log.Println("Error acquiring payment lock:", err)
		return
	}
	log.Println("Lock acquired for order:", event.OrderID)
	time.Sleep(5 * time.Second)
	if !acquired {
		log.Println("Payment already being processed for order:", event.OrderID)
		return
	}
	defer func() {
		if err := c.lock.Release(ctx, event.OrderID); err != nil {
			log.Println("Error releasing payment lock:", err)
		}
	}()
	isCompleted, price, err := c.repo.ProcessPayment(
		event.OrderID,
		float32(event.TotalPrice),
		event.CustomerID,
	)
	if err != nil {
		log.Println("Payment error:", err)
		return
	}

	if !isCompleted {
		log.Println("Payment failed for order:", event.OrderID)
		return
	}
	err = c.producer.PublishEvent(
		event.OrderID,
		event.CustomerID,
		"COMPLETED",
		price,
	)
	if err != nil {
		log.Println("Error publishing payment.completed event:", err)
		return
	}

	log.Println(
		"Payment completed for order:",
		event.OrderID,
		"price:",
		price,
	)
}
