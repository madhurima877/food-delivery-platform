package main

import (
	"log"
	"net"

	"github.com/madhurima877/food-delivery-platform/proto/notification"
	"github.com/madhurima877/food-delivery-platform/services/notification-service/handler"
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
	log.Println("Notification Service started on port 50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
