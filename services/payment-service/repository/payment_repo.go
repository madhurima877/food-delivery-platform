package repository

import (
	"database/sql"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}
func (repo *PaymentRepository) ProcessPayment(orderID string, price float32, userID string) (string, float32, error) {

	query := `
		INSERT INTO payments (order_id, user_id, amount, status)
VALUES ($1, $2, $3, $4)
ON CONFLICT (order_id) DO NOTHING
	`

	result, err := repo.db.Exec(
		query,
		orderID,
		userID,
		price,
		"COMPLETED",
	)
	if err != nil {
		return "FAILED", 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "FAILED", 0, err
	}

	if rowsAffected == 0 {
		return "DUPLICATE", price, nil
	}

	return "COMPLETED", price, nil
}
