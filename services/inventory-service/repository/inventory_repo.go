package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (repo *InventoryRepository) ReserveStock(orderID string, productID string, quantity int32) (string, int32, error) {
	if productID == "999" {
		return "", 0, fmt.Errorf("simulated database error")
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return "", 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO inventory_processed_events (order_id)
		VALUES ($1)
		ON CONFLICT (order_id) DO NOTHING
	`, orderID)

	if err != nil {
		return "", 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", 0, err
	}

	if rowsAffected == 0 {
		log.Println("Duplicate order skipped:", orderID)
		return "DUPLICATE", 0, nil
	}

	var leftStock int32

	query := `
		UPDATE inventory
		SET
			available_stock = available_stock - $1,
			reserved_stock = reserved_stock + $1
		WHERE product_id = $2
		AND available_stock >= $1
		RETURNING available_stock
	`

	err = tx.QueryRow(query, quantity, productID).Scan(&leftStock)

	if err == sql.ErrNoRows {
		return "NOT_ENOUGH_STOCK", 0, nil
	}

	if err != nil {
		return "", 0, err
	}

	if err := tx.Commit(); err != nil {
		return "", 0, err
	}

	return "RESERVED", leftStock, nil
}

func (repo *InventoryRepository) GetProductPrice(productID string) (int, error) {
	var price int

	query := `SELECT price FROM products WHERE id = $1`

	err := repo.db.QueryRow(query, productID).Scan(&price)
	if err != nil {
		return 0, err
	}

	return price, nil
}
