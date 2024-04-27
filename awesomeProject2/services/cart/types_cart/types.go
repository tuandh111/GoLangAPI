package types_cart

import "time"

type CartCheckoutItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}
type CartItem struct {
	Id        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreateAt  time.Time `json:"createAt"`
}
type CartItemUpdate struct {
	OrderID   int     `json:"orderID"`
	ProductID int     `json:"productID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CartCheckoutPayload struct {
	Items []CartCheckoutItem `json:"items" validate:"required"`
}
type CartStore interface {
	CartItems(orderId int) ([]*CartItem, error)
	CheckCartOrderIdAndProductId(orderId int, productId int) (*CartItem, int)
	UpdateOrSaveOrderIdAndProductId(cartItemUpdate CartItemUpdate, cartID int) (string, error)
}
