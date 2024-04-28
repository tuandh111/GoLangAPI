package types_checkout

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}
type OrderCheckout struct {
	ID     int    `json:"id" validation:"required"`
	Status string `json:"status"  validation:"required"`
}

type CheckOutStore interface {
	UpdateStatusAdmin(orderItemId int, status string) (string, error)
}
