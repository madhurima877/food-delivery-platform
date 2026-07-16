package repository

import "database/sql"

type OrderRepository struct {
	db *sql.DB
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
