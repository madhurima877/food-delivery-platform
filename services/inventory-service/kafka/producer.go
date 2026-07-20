package kafka

import (
	"context"
	"encoding/json"
	"strconv"

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
	return &Producer{
		writer: writer,
	}
}
func (p *Producer) Close() error {
	return p.writer.Close()
}
func (p *Producer) PublishEvent(OrderID, CustomerID, ProductID string, TotalPrice, Quantity int32) error {
	msg := models.InventoryReservedEvent{
		OrderID:    OrderID,
		CustomerID: CustomerID,
		ProductID:  ProductID,
		Quantity:   Quantity,
		TotalPrice: TotalPrice,
	}
	msgbyte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: "inventory.reserved",
			Value: msgbyte,
		},
	)
	if err != nil {
		return err
	}
	return nil

}
func (p *Producer) PublishRetryEvent(message []byte, retryCount int) error {
	msg := kafka.Message{
		Value: message,
		Topic: "order.created.retry",
		Headers: []kafka.Header{
			{
				Key:   "retry-count",
				Value: []byte(strconv.Itoa(retryCount)),
			},
		},
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

func (p *Producer) PublishDLQEvent(message []byte) error {
	msg := kafka.Message{
		Value: message,
		Topic: "order.created.dlq",
	}
	return p.writer.WriteMessages(context.Background(), msg)
}
