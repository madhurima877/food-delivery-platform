package repository

import (
	"database/sql"
)

type OrderRepository struct {
	db *sql.DB
}

type Order struct {
	Id           string
	CustomerId   string
	RestaurantId string
	Status       string
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (repo *OrderRepository) CreateOrder(customerID, restaurantID string) (int64, error) {
	var id int64
	sql := `INSERT INTO orders (customer_id,restaurant_id,status)VALUES($1,$2,$3) RETURNING id`

	err := repo.db.QueryRow(sql, customerID, restaurantID, "CREATED").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *OrderRepository) GetOrder(orderId string) (*Order, error) {
	var order Order

	sql := `
		SELECT id, customer_id, restaurant_id, status
		FROM orders
		WHERE id = $1
	`

	err := repo.db.QueryRow(sql, orderId).Scan(
		&order.Id,
		&order.CustomerId,
		&order.RestaurantId,
		&order.Status,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (repo *OrderRepository) UpdateOrder(orderId, status string) (*Order, error) {
	var order Order
	sql := `Update orders SET status=$1 WHERE id=$2 RETURNING id, customer_id, restaurant_id, status`
	err := repo.db.QueryRow(sql, orderId, status).Scan(order.Id, order.CustomerId, order.RestaurantId, order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
