package main

import (
	"log"
	"net"

	"github.com/madhurima877/food-delivery-platform/proto/payment"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/db"
	"github.com/madhurima877/food-delivery-platform/services/payment-service/handler"
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
	log.Println("Payment Service started on port 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
