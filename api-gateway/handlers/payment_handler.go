package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	pbP "github.com/madhurima877/food-delivery-platform/proto/payment"
)

func ProcessPaymentHandler(client pbP.PaymentServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.ProcessPaymentRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.ProcessPayment(context.Background(), &pbP.ProcessPaymentRequest{
			Price:   req.Price,
			OrderId: req.OrderId,
			UserId:  req.UserId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}

}
