package models

type CreateOrderRequest struct {
	CustomerID   string `json:"customer_id"`
	RestaurantID string `json:"restaurant_id"`
	ProductID    string `json:"product_id"`
	Quantity     int32  `json:"quantity"`
}
type UpdateOrderRequest struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}
type InventoryReserveRequest struct {
	OrderId   string `json:"order_id"`
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}
type ProcessPaymentRequest struct {
	OrderId string  `json:"order_id"`
	Price   float32 `json:"price"`
	UserId  string  `json:"user_id"`
}
type NotificationRequest struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

type DriverRequest struct {
	OrderId  string `json:"order_id"`
	DriverId string `json:"driver_id"`
}
type OrderCreatedEvent struct {
	OrderID      string `json:"order_id"`
	CustomerID   string `json:"customer_id"`
	RestaurantID string `json:"restaurant_id"`
	ProductID    string `json:"product_id"`
	Quantity     int32  `json:"quantity"`
}
type InventoryReservedEvent struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	ProductID  string `json:"product_id"`
	Quantity   int32  `json:"quantity"`
	TotalPrice int32  `json:"total_price"`
}
type PaymentCompletedEvent struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	Amount     float32 `json:"amount"`
	Status     string  `json:"status"`
}
type Order struct {
	Id           string
	CustomerId   string
	RestaurantId string
	Status       string
}
type PaymentFailedEvent struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	ProductID  string `json:"product_id"`
	Quantity   int32  `json:"quantity"`
}
