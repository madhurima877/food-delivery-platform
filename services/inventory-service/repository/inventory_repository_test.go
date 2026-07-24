package repository

import (
	"testing"

	"github.com/madhurima877/food-delivery-platform/services/notification-service/db"
)

func TestRestoreStock(t *testing.T) {
	database, err := db.Connect()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryRepository(database)
	_, err = database.Exec(`
	INSERT INTO inventory (product_id, available_stock, reserved_stock)
	VALUES ($1, $2, $3)
	ON CONFLICT (product_id) DO UPDATE
	SET available_stock = $2, reserved_stock = $3
`, 999, 90, 10)

	if err != nil {
		t.Fatal(err)
	}
	err = repo.RestoreStock("999", 2, "test-order-1")
	if err != nil {
		t.Fatal(err)
	}
	var availableStock int
	var reservedStock int
	err = database.QueryRow(`SELECT available_stock, reserved_stock
	FROM inventory
	WHERE product_id = $1`, 999).Scan(&availableStock, &reservedStock)

	if err != nil {
		t.Fatal(err)
	}
	if availableStock != 92 {
		t.Errorf("expected available stock 92, got %d", availableStock)
	}

	if reservedStock != 8 {
		t.Errorf("expected reserved stock 8, got %d", reservedStock)
	}
}
