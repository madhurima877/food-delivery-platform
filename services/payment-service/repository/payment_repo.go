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
func (repo *PaymentRepository) ProcessPayment(orderID string, price float32, userID string) (bool, float32, error) {
	query := `
		INSERT INTO payments (order_id, user_id, amount, status)
		VALUES ($1, $2, $3, $4)
	`

	result, err := repo.db.Exec(
		query,
		orderID,
		userID,
		price,
		"COMPLETED",
	)
	if err != nil {
		return false, 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, 0, err
	}

	if rowsAffected == 0 {
		return false, 0, nil
	}

	return true, price, nil
}
