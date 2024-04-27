package cart

import (
	"awesomeProject2/services/cart/types_cart"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewStoreCart(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
func (s *Store) CartItems(orderId int) ([]*types_cart.CartItem, error) {
	rows, err := s.db.Query("select * from order_items where orderID = ?", orderId)
	if err != nil {
		return nil, err
	}
	order := make([]*types_cart.CartItem, 0)
	for rows.Next() {
		cart, errs := scanRowsIntoOrderItem(rows)
		if errs != nil {
			return nil, errs
		}
		order = append(order, cart)
	}
	return order, nil
}
func (s *Store) CheckCartOrderIdAndProductId(orderId int, productId int) (*types_cart.CartItem, int) {
	row := s.db.QueryRow("select * from order_items where orderID = ? and productID = ?", orderId, productId)
	order := new(types_cart.CartItem)
	cart, err := scanRowIntoOrderItem(row, order)
	if err != nil {
		return nil, 0
	}
	return cart, cart.Id
}
func (s *Store) UpdateOrSaveOrderIdAndProductId(cartItemUpdate types_cart.CartItemUpdate, OrderItemId int) (string, error) {
	_, err := s.db.Exec("update order_items set orderID = ?, productID = ?, quantity = ?, price = ? where id = ?", cartItemUpdate.OrderID, cartItemUpdate.ProductID, cartItemUpdate.Quantity, cartItemUpdate.Price, OrderItemId)
	if err != nil {
		return "", err
	}

	return "update successfully", nil
}
func scanRowsIntoOrderItem(rows *sql.Rows) (*types_cart.CartItem, error) {
	orderItem := new(types_cart.CartItem)
	if err := rows.Scan(&orderItem.Id, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price, &orderItem.CreateAt); err != nil {
		return nil, err
	}
	return orderItem, nil
}

func scanRowIntoOrderItem(rows *sql.Row, orderItem *types_cart.CartItem) (*types_cart.CartItem, error) {
	if err := rows.Scan(&orderItem.Id, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price, &orderItem.CreateAt); err != nil {
		return nil, err
	}
	return orderItem, nil
}
