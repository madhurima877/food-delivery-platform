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

	sql := `
		UPDATE inventory
		SET
			available_stock = available_stock - $1,
			reserved_stock = reserved_stock + $1
		WHERE product_id = $2
		AND available_stock >= $1
		RETURNING available_stock
	`

	err := repo.db.QueryRow(sql, quantity, productID).Scan(&leftStock)
	if err != nil {
		return false, 0, err
	}

	return true, leftStock, nil
}
