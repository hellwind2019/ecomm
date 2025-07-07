package storer

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MySQLStorer struct {
	db *sqlx.DB
}

func NewMySQLStorer(db *sqlx.DB) *MySQLStorer {
	return &MySQLStorer{
		db: db,
	}
}

func (s *MySQLStorer) CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	query := `INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)`

	res, err := s.db.NamedExecContext(ctx, query, p)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	p.ID = id
	return p, nil
}

func (s *MySQLStorer) GetProduct(ctx context.Context, id int64) (*Product, error) {
	var p Product
	err := s.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &p, nil
}
func (s *MySQLStorer) ListProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	err := s.db.SelectContext(ctx, &products, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

func (s *MySQLStorer) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	query := `
		UPDATE products SET
			name = :name,
			image = :image,
			category = :category,
			description = :description,
			rating = :rating,
			num_reviews = :num_reviews,
			price = :price,
			count_in_stock = :count_in_stock,
			updated_at = NOW()
		WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, query, p)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return p, nil
}
func (s *MySQLStorer) DeleteProduct(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
