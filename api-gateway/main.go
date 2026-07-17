package main

import (
	"fmt"
	"net/http"

	"github.com/madhurima877/food-delivery-platform/api-gateway/handlers"

	pbb "github.com/madhurima877/food-delivery-platform/proto/inventory"
	pb "github.com/madhurima877/food-delivery-platform/proto/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	orderconn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	//order
	orderclient := pb.NewOrderServiceClient(orderconn)

	http.HandleFunc("/orders", handlers.CreateOrderHandler(orderclient))
	http.HandleFunc("/get/order", handlers.GetOrderHandler(orderclient))
	http.HandleFunc("/update/status", handlers.UpdateOrderHandler(orderclient))
	http.HandleFunc("/delete/order", handlers.DeleteOrderHandler(orderclient))

	//inventory service

	invetoryconn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	inventoryClient := pbb.NewInventoryServiceClient(invetoryconn)

	http.HandleFunc(
		"/inventory/reserve",
		handlers.ReserveStockHandler(inventoryClient),
	)

	fmt.Println("API Gateway started on :8080")

	http.ListenAndServe(":8080", nil)
}
