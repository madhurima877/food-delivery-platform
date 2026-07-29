package repository

import (
	"testing"

	"github.com/madhurima877/food-delivery-platform/services/payment-service/db"
)

func TestProcessPaymentIdempotency(t *testing.T) {
	database, err := db.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repo := NewPaymentRepository(database)
	_, err = database.Exec(`DELETE FROM payments
	WHERE order_id = $1
`, 999)

	if err != nil {
		t.Fatal(err)
	}
	status, price, err := repo.ProcessPayment("999", 250, "101")
	if err != nil {
		t.Fatal(err)
	}
	if status != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", status)
	}
	if price != 250 {
		t.Errorf("expected price 250, got %.2f", price)
	}
	status, price, err = repo.ProcessPayment("999", 250, "101")
	if err != nil {
		t.Fatal(err)
	}
	if status != "DUPLICATE" {
		t.Errorf("expected DUPLICATE, got %s", status)
	}
	var count int
	err = database.QueryRow(`
	SELECT COUNT(*) FROM payments WHERE order_id = $1
`, 999).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected count to be 1, got %d", count)
	}
}
