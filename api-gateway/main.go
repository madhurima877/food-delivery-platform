package main

import (
	"fmt"
	"net/http"

	"github.com/madhurima877/food-delivery-platform/api-gateway/handlers"
	"github.com/madhurima877/food-delivery-platform/api-gateway/ratelimit"
	pbD "github.com/madhurima877/food-delivery-platform/proto/driver"
	pbN "github.com/madhurima877/food-delivery-platform/proto/notification"
	pbP "github.com/madhurima877/food-delivery-platform/proto/payment"

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

	ratelimiter := ratelimit.NewRateLimit()

	//order
	orderclient := pb.NewOrderServiceClient(orderconn)

	http.HandleFunc("/orders", ratelimiter.Middleware(handlers.CreateOrderHandler(orderclient)))
	http.HandleFunc("/get/order", ratelimiter.Middleware(handlers.GetOrderHandler(orderclient)))
	http.HandleFunc("/update/status", ratelimiter.Middleware(handlers.UpdateOrderHandler(orderclient)))
	http.HandleFunc("/delete/order", ratelimiter.Middleware(handlers.DeleteOrderHandler(orderclient)))

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

	paymentconn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	paymentClient := pbP.NewPaymentServiceClient(paymentconn)

	http.HandleFunc("/payment", handlers.ProcessPaymentHandler(paymentClient))

	notificationconn, err := grpc.NewClient("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	notificationClient := pbN.NewNotificationServiceClient(notificationconn)
	http.HandleFunc("/send/notification", handlers.SendNotificationHandler(notificationClient))

	driverconn, err := grpc.NewClient("localhost:50056", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	driverClient := pbD.NewDriverServiceClient(driverconn)
	http.HandleFunc("/assign/driver", handlers.AssignDriverHandler(driverClient))

	fmt.Println("API Gateway started on :8080")

	http.ListenAndServe(":8080", nil)
}
