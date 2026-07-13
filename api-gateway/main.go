package main

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/madhurima877/food-delivery-platform/proto/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createOrderHandler(client pb.OrderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resp, err := client.CreateOrder(context.Background(), &pb.CreateOrderRequest{
			CustomerId:   "123",
			RestaurantId: "456",
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, resp.OrderId, resp.Status)
	}

}
func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := pb.NewOrderServiceClient(conn)

	http.HandleFunc("/orders", createOrderHandler(client))
	fmt.Println("API Gateway started on :8080")

	http.ListenAndServe(":8080", nil)
}
