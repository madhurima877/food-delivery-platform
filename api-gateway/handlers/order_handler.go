package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	pb "github.com/madhurima877/food-delivery-platform/proto/order"
)

func CreateOrderHandler(client pb.OrderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateOrderRequest
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

func UpdateOrderHandler(client pb.OrderServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UpdateOrderRequest
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
