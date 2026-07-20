package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/madhurima877/food-delivery-platform/proto/driver"
	"github.com/madhurima877/food-delivery-platform/services/driver-service/handler"
	"github.com/madhurima877/food-delivery-platform/services/driver-service/kafka"
	"github.com/madhurima877/food-delivery-platform/services/driver-service/repository"
	"google.golang.org/grpc"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered From Panic")
		}
	}()
	lis, err := net.Listen("tcp", ":50056")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	repo := repository.NewDriverRepository()
	driverhandler := handler.NewDriverHandler(repo)
	driver.RegisterDriverServiceServer(grpcServer, driverhandler)
	go func() {
		log.Println("Driver Service started on port 50056")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewConsumer(repo)
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go consumer.ReadConsumer(ctx, &wg)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down Driver Service...")

	cancel()
	wg.Wait()
	grpcServer.GracefulStop()

	log.Println("Driver Service stopped")

}
