package product

import (
	"awesomeProject2/types"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
func (s *Store) GetProductByID(id int) (*types.Product, error) {

	row := s.db.QueryRow("select  * from products where id = ?", id)
	p := &types.Product{}
	err := scanRowIntoProduct(row, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (s *Store) GetProductsByID(ids []int) ([]types.Product, error) {
	var products []types.Product
	return products, nil
}
func (s *Store) GetProductsPage(page int, limit int) ([]*types.Product, error) {
	products := make([]*types.Product, 0)
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT * FROM products LIMIT %d OFFSET %d", limit, offset)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p, errs := scanRowsIntoProduct(rows)
		if errs != nil {
			return nil, errs
		}
		products = append(products, p)
	}
	return products, nil
}
func (s *Store) GetProducts() ([]*types.Product, error) {
	products := make([]*types.Product, 0)
	rows, err := s.db.Query("select  * from products")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p, errs := scanRowsIntoProduct(rows)
		if errs != nil {
			return nil, errs
		}
		products = append(products, p)
	}
	return products, nil
}
func (s *Store) CreateProduct(types.CreateProductPayload) error {
	return nil
}
func (s *Store) UpdateProduct(types.Product) error {
	return nil
}
func (s *Store) DeleteProductByID(id int) error {
	return nil
}
func scanRowIntoProduct(row *sql.Row, p *types.Product) error {
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	user := new(types.Product)
	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Description,
		&user.Image,
		&user.Price,
		&user.Quantity,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
