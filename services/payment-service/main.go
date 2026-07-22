package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/madhurima877/food-delivery-platform/proto/payment"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/db"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/handler"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/kafka"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/lock"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/repository"
	"google.golang.org/grpc"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		panic(err)
	}
	log.Println("Database Connected")

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	repo := repository.NewPaymentRepository(database)
	paymentHandler := handler.NewPaymentHandler(repo)
	payment.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	go func() {
		log.Println("Payment Service started on port 50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	producer := kafka.NewProducer()
	defer producer.Close()
	redisLock := lock.NewRedisLock()

	consumer := kafka.NewConsumer(repo, producer, redisLock)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go consumer.ReadConsumer(ctx, &wg)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	log.Println("Shutting down Payment Service...")

	cancel()                  // 1. Tell Kafka workers to stop
	grpcServer.GracefulStop() // 2. Stop gRPC server
	wg.Wait()                 // 3. Wait for all 3 workers to finish

	log.Println("Payment Service stopped")

}
