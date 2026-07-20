package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/madhurima877/food-delivery-platform/proto/notification"
	"github.com/madhurima877/food-delivery-platform/services/notification-service/handler"
	"github.com/madhurima877/food-delivery-platform/services/notification-service/kafka"
	"github.com/madhurima877/food-delivery-platform/services/notification-service/repository"
	"google.golang.org/grpc"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered From Panic")
		}
	}()

	// database, err := db.Connect()
	// if err != nil {
	// 	panic(err)
	// }
	log.Println("Database Connected")
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	repo := repository.NewNotificationRepository()
	notificationHandler := handler.NewNotificationHandler(repo)
	notification.RegisterNotificationServiceServer(grpcServer, notificationHandler)
	go func() {
		log.Println("Notification Service started on port 50054")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	consumer := kafka.NewConsumer(repo)

	log.Println("Starting Notification Kafka consumer workers...")

	for i := 1; i <= 3; i++ {
		wg.Add(1)

		go func(workerID int) {
			log.Println("Notification Kafka worker started:", workerID)
			consumer.ReadConsumer(ctx, &wg)
		}(i)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	log.Println("Shutting down Notification Service...")

	cancel()
	wg.Wait()
	grpcServer.GracefulStop()

	log.Println("Notification Service stopped")

}
