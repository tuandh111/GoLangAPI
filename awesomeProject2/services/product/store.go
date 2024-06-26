package product

import (
	"awesomeProject2/services/product/types_product"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
func (s *Store) GetProductByID(id int) (*types_product.Product, error) {

	row := s.db.QueryRow("select  * from products where id = ?", id)
	p := &types_product.Product{}
	err := scanRowIntoProduct(row, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (s *Store) GetProductsByID(ids []int) ([]types_product.Product, error) {
	var products []types_product.Product
	return products, nil
}
func (s *Store) GetProductsPage(page int, limit int) ([]*types_product.Product, error) {
	products := make([]*types_product.Product, 0)
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
func (s *Store) GetProducts() ([]*types_product.Product, error) {
	products := make([]*types_product.Product, 0)
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
func (s *Store) GetProductByName(productName string) (*types_product.Product, error) {
	rows := s.db.QueryRow("select  * from products where name = ? ", productName)
	p := &types_product.Product{}
	err := scanRowIntoProduct(rows, p)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return p, nil

}
func (s *Store) CreateProduct(product types_product.CreateProductPayload) error {
	_, err := s.db.Exec("insert into products ( name,description, image, price, quantity) values (?,?,?,?,?)", product.Name, product.Description, product.Image, product.Price, product.Quantity)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) UpdateProduct(product types_product.UpdateProduct, productId int) (string, error) {
	_, err := s.db.Exec("update products set name = ?, description = ?, image = ? , price = ? , quantity = ?  where id = ?", product.Name, product.Description, product.Image, product.Price, product.Quantity, productId)
	if err != nil {
		return "error update ", err
	}
	return "update successfully", nil
}
func (s *Store) DeleteProductByID(id int) error {
	return nil
}
func scanRowIntoProduct(row *sql.Row, p *types_product.Product) error {
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func scanRowsIntoProduct(rows *sql.Rows) (*types_product.Product, error) {
	user := new(types_product.Product)
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
