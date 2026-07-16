package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/madhurima877/food-delivery-platform/proto/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateOrderRequest struct {
	CustomerID   string `json:"customer_id"`
	RestaurantID string `json:"restaurant_id"`
}

func createOrderHandler(client pb.OrderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateOrderRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp, err := client.CreateOrder(context.Background(), &pb.CreateOrderRequest{
			CustomerId:   req.CustomerID,
			RestaurantId: req.RestaurantID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, resp.OrderId, resp.Status)
	}

}

type UpdateOrderRequest struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}

func UpdateOrderHandler(client pb.OrderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateOrderRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp, err := client.UpdateOrderStatus(context.Background(), &pb.UpdateOrderStatusRequest{
			OrderId: req.OrderId,
			Status:  req.Status,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}
func GetOrderHandler(client pb.OrderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := r.URL.Query().Get("order_id")
		resp, err := client.GetOrder(context.Background(), &pb.GetOrderRequest{OrderId: orderId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)

	}

}

func DeleteOrderHandler(client pb.OrderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderid := r.URL.Query().Get("order_id")
		resp, err := client.DeleteOrder(context.Background(), &pb.DeleteOrderRequest{OrderId: orderid})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}
func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := pb.NewOrderServiceClient(conn)

	http.HandleFunc("/orders", createOrderHandler(client))
	http.HandleFunc("/get/order", GetOrderHandler(client))
	http.HandleFunc("/update/status", UpdateOrderHandler(client))
	http.HandleFunc("/delete/order", DeleteOrderHandler(client))
	fmt.Println("API Gateway started on :8080")

	http.ListenAndServe(":8080", nil)
}
