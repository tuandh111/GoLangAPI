package order

import (
	"awesomeProject2/services/order/types_order"
	"awesomeProject2/types"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewOrder(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
func (s *Store) FindAllOrderWithAdmin() ([]*types.Order, error) {
	rows, err := s.db.Query("select  * from orders")
	if err != nil {
		return nil, err
	}
	orders := make([]*types.Order, 0)
	for rows.Next() {
		row, err := scanRowsIntoOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, row)
	}
	return orders, nil
}

func scanRowsIntoOrder(row *sql.Rows) (*types.Order, error) {
	order := new(types.Order)
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.Address,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (s *Store) CreateOrder(order types_order.OrderPayload) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserId, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}
	id, er := res.LastInsertId()
	if er != nil {
		return 0, er
	}

	return int(id), nil
}
func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return err
	}
	return nil
}
