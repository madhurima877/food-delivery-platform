package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	pbb "github.com/madhurima877/food-delivery-platform/proto/inventory"
)

func ReserveStockHandler(client pbb.InventoryServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.InventoryReserveRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.ReserveStock(
			context.Background(),
			&pbb.ReserveStockRequest{
				OrderId:   req.OrderId,
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}

}
