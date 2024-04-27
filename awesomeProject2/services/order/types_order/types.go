package types_order

import "time"

type OrderPayload struct {
	UserId  string `json:"userID" `
	Total   int    `json:"total" validate:"required"`
	Status  string `json:"status" validate:"required"`
	Address string `json:"address" validate:"required"`
}
type OrderUpdateUserID struct {
	Total   int    `json:"total" validate:"required"`
	Status  string `json:"status" validate:"required"`
	Address string `json:"address" validate:"required"`
}
type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}
type OrderStore interface {
	CreateOrder(payload OrderPayload) (int, error)
	CreateOrderItem(OrderItem) error
	FindAllOrderWithAdmin() ([]*Order, error)
	FindByOrderUserId(userId int) ([]*Order, error)
	FindByOrderUserIdAndStatus(userId int, status string) (*Order, error)
	//PENDING,PROCESSING,SHIPPING,DELIVERED,CANCELLED
	UpdateOrderByUserId(orderUpdate OrderUpdateUserID, userId int, status string) (string, error)
}
