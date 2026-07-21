package main

import (
	"context"
	"log"
	"net"

	"github.com/madhurima877/food-delivery-platform/proto/order"
	"github.com/madhurima877/food-delivery-platform/services/order-service/cache"
	"github.com/madhurima877/food-delivery-platform/services/order-service/db"
	"github.com/madhurima877/food-delivery-platform/services/order-service/handler"
	"github.com/madhurima877/food-delivery-platform/services/order-service/kafka"
	"github.com/madhurima877/food-delivery-platform/services/order-service/repository"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	redisCache := cache.NewRedisCache()
	if err := redisCache.Ping(ctx); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis Connected")
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

	orderHandler := handler.NewOrderHandler(repo, writer, redisCache)

	order.RegisterOrderServiceServer(grpcServer, orderHandler)

	log.Println("Order Service started on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
