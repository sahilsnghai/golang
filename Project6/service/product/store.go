package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sahilsnghai/Project6/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("Select * from products")

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

func (s *Store) GetProductByID(productId int) (*types.Product, error) {
	rows, err := s.db.Query("select * from products where id = ?", productId)
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
	return p, nil
}



func (s *Store) GetProductsByID(productIds []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIds)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	// Convert productIds to []interface{}
	args := make([]interface{}, len(productIds))
	for i, v := range productIds {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil

}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products set name = ?, price = ?, image = ?, description = ?, quantity = ? where id =?",
		product.Name, product.Price, product.Image, product.Description, product.Quantity, product.Id)

	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}
