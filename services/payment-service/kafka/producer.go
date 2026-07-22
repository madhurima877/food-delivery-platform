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
		Addr: kafka.TCP("localhost:19092"),
	}
	return &Producer{writer: writer}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) PublishEvent(OrderID, CustomerID, Status string, Amount float32) error {
	event := models.PaymentCompletedEvent{
		OrderID:    OrderID,
		CustomerID: CustomerID,
		Status:     Status,
		Amount:     Amount,
	}

	eventByte, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: "payment.completed",
			Value: eventByte,
		})
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) PublishFailedEvent(orderID string, customerID string, productID string, quantity int32) error {
	event := models.PaymentFailedEvent{
		OrderID:    orderID,
		CustomerID: customerID,
		ProductID:  productID,
		Quantity:   quantity,
	}
	eventByte, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Value: eventByte,
		Topic: "payment.failed",
	})

}
