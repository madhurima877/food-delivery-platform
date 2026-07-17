package main

import (
	"log"
	"net"

	"github.com/madhurima877/food-delivery-platform/proto/inventory"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/db"
	"github.com/madhurima877/food-delivery-platform/services/inventory-service/handler"
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
	log.Println("Inventory Service started on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
