package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/madhurima877/food-delivery-platform/api-gateway/models"
	pbN "github.com/madhurima877/food-delivery-platform/proto/notification"
)

func SendNotificationHandler(client pbN.NotificationServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.NotificationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.SendNotification(context.Background(), &pbN.NotificationRequest{
			UserId:  req.UserId,
			Message: req.Message,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}
