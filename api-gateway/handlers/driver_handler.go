package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	pbD "github.com/madhurima877/food-delivery-platform/proto/driver"
)

func AssignDriverHandler(client pbD.DriverServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.DriverRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.AssignDriver(context.Background(), &pbD.DriverRequest{
			OrderId:  req.OrderId,
			DriverId: req.DriverId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)

	}
}
