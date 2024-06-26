package order

import (
	"awesomeProject2/services/order/types_order"
	"database/sql"
	"fmt"
	"strconv"
)

type Store struct {
	db *sql.DB
}

func NewOrder(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
func (s *Store) FindAllOrderWithAdmin() ([]*types_order.Order, error) {
	rows, err := s.db.Query("select  * from orders")
	if err != nil {
		return nil, err
	}
	orders := make([]*types_order.Order, 0)
	for rows.Next() {
		row, errs := scanRowsIntoOrder(rows)
		if errs != nil {
			return nil, err
		}
		orders = append(orders, row)
	}
	return orders, nil
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
func (s *Store) CreateOrderItem(orderItem types_order.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) FindByOrderUserId(userId int) ([]*types_order.Order, error) {
	rows, err := s.db.Query("select  * from orders where userID = ?", userId)
	if err != nil {
		return nil, err
	}
	order := make([]*types_order.Order, 0)
	for rows.Next() {
		o, err := scanRowsIntoOrder(rows)
		if err != nil {
			return nil, err
		}
		order = append(order, o)

	}
	return order, nil
}
func (s *Store) UpdateOrderByUserId(payload types_order.OrderUpdateUserID, userId int, status string) (string, error) {
	_, err := s.db.Exec("update orders set total = ?, status = ?, address = ? where userID = ? and status = ?", payload.Total, payload.Status, payload.Address, userId, status)
	if err != nil {
		return "fail update order", err
	}
	return "update successfully", nil
}
func (s *Store) FindByOrderUserIdAndStatus(userId int, status string) (*types_order.Order, error) {
	row := s.db.QueryRow("select * from orders where userID = ? and status = ?", userId, status)
	order := &types_order.Order{}
	if err := scanRowIntoOrder(row, order); err != nil {
		return nil, err
	}
	return order, nil
}
func (s *Store) DeleteOrder(orderId int) (string, error) {
	_, errDeleteCart := s.db.Exec("delete from order_items where orderID = ?", orderId)
	if errDeleteCart != nil {
		fmt.Println(errDeleteCart)
	}
	_, err := s.db.Exec("delete from orders where id = ?", orderId)
	if err != nil {
		return "delete error: ", err
	}
	return "delete successfully with orderId: " + strconv.Itoa(orderId), nil
}

func scanRowIntoOrder(row *sql.Row, order *types_order.Order) error {
	if err := row.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.Address, &order.CreatedAt); err != nil {
		return err
	}
	return nil
}
func scanRowsIntoOrder(row *sql.Rows) (*types_order.Order, error) {
	order := new(types_order.Order)
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
