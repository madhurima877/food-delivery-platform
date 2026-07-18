package repository

import "log"

type DriverRepository struct {
}

func NewDriverRepository() *DriverRepository {
	return &DriverRepository{}
}
func (repo *DriverRepository) AssignDriver(orderId, driverID string) {
	log.Println("DriverId ", driverID, "assigned to OrderId ", orderId)

}
