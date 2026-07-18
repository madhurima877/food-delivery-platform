package main

import (
	"log"
	"net"

	"github.com/madhurima877/food-delivery-platform/proto/driver"
	"github.com/madhurima877/food-delivery-platform/services/driver-service/handler"
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

	log.Println("Driver Service started on port 50056")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
