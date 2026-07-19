package repository

import "database/sql"

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (repo *InventoryRepository) ReserveStock(productID string, quantity int32) (bool, int32, error) {
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

	err := repo.db.QueryRow(query, quantity, productID).Scan(&leftStock)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	return true, leftStock, nil
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
