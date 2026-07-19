package main

import (
	"log"
	"net"

	"github.com/madhurima877/food-delivery-platform/proto/order"
	"github.com/madhurima877/food-delivery-platform/services/order-service/db"
	"github.com/madhurima877/food-delivery-platform/services/order-service/handler"
	"github.com/madhurima877/food-delivery-platform/services/order-service/kafka"
	"github.com/madhurima877/food-delivery-platform/services/order-service/repository"

	"google.golang.org/grpc"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	log.Println("Database Connected")

	writer := kafka.NewProducer()
	defer writer.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	repo := repository.NewOrderRepository(database)

	orderHandler := handler.NewOrderHandler(repo, writer)

	order.RegisterOrderServiceServer(grpcServer, orderHandler)

	log.Println("Order Service started on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
