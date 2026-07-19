package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/madhurima877/food-delivery-platform/proto/inventory"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/db"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/handler"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/kafka"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/repository"
	"google.golang.org/grpc"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Connected!!!!")
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	inventoryRepo := repository.NewInventoryRepository(database)

	inventoryHandler := handler.NewInventoryHandler(inventoryRepo)

	inventory.RegisterInventoryServiceServer(grpcServer, inventoryHandler)
	// log.Println("Inventory Service started on port 50052")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	producer := kafka.NewProducer()
	defer producer.Close()
	consumer := kafka.NewConsumer(inventoryRepo, producer)

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go consumer.ReadConsumer(ctx, &wg)
	}
	go func() {
		log.Println("Inventory Service started on port 50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("grpc server stopped", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down Inventory Service...")

	cancel()                  // tells Kafka workers to stop
	grpcServer.GracefulStop() // stops gRPC
	wg.Wait()                 // waits for all 4 workers

	log.Println("Inventory Service stopped")

}
