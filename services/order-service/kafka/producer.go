package kafka

import (
	"context"
	"encoding/json"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:19092"),
		Topic: "order.created",
	}
	return &Producer{writer: writer}
}
func (p *Producer) Close() error {
	return p.writer.Close()
}
func (p *Producer) PublishOrderCreated(orderID string,
	customerID string,
	restaurantID string,
	productID string,
	quantity int32) error {
	msg := models.OrderCreatedEvent{
		OrderID:      orderID,
		CustomerID:   customerID,
		RestaurantID: restaurantID,
		ProductID:    productID,
		Quantity:     quantity,
	}

	msgbyte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = p.writer.WriteMessages(context.Background(), kafka.Message{Value: msgbyte})
	if err != nil {
		return err
	}
	return nil
}
