package product

import (
	"database/sql"
	"fmt"

	"github.com/anangfirmansyahp5/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, price, quantity) VALUES (?, ?, ?, ?)", product.Name, product.Description, product.Price, product.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProductByID(id int) (*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	p := new(types.Product)
	for rows.Next() {
		p, err = scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return p, nil
}

func (s *Store) UpdateProduct(id int, product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, price = ?, quantity = ? WHERE id = ?", product.Name, product.Description, product.Price, product.Quantity, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteProduct(id int) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}
