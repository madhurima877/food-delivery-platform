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
		Topic: "inventory.reserved",
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
	err = p.writer.WriteMessages(context.Background(), kafka.Message{Value: msgbyte})
	if err != nil {
		return err
	}
	return nil

}
