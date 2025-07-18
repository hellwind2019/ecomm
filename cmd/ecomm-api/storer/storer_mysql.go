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
func (s *MySQLStorer) ListProducts(ctx context.Context) ([]Product, error) {
	var products []Product
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
			updated_at = :updated_at
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
func (s *MySQLStorer) CreateOrder(ctx context.Context, o *Order) (*Order, error) {
	err := s.execTx(ctx, func(tx *sqlx.Tx) error {
		//insert into orders
		order, err := createOrder(ctx, tx, o)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
		for _, oi := range o.Items {
			oi.OrderID = order.ID // Set the OrderID for each OrderItem
			err = createOrderItems(ctx, tx, oi)
			if err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}
		return nil
		//insert into order_items
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return o, nil
	//start a transaction

	//commit the transaction
	//rollback the transaction if any error occurs

}
func createOrder(ctx context.Context, tx *sqlx.Tx, o *Order) (*Order, error) {
	res, err := tx.NamedExecContext(ctx, `
        INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, user_id)
        VALUES (:payment_method, :tax_price, :shipping_price, :total_price, :user_id)
    `, o)
	if err != nil {
		return nil, fmt.Errorf("error inserting order: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", err)
	}
	o.ID = id
	return o, nil
}
func createOrderItems(ctx context.Context, tx *sqlx.Tx, oi OrderItem) error {
	res, err := tx.NamedExecContext(ctx, `INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (:name, :quantity, :image, :price, :product_id, :order_id)`, oi)
	if err != nil {
		return fmt.Errorf("error inserting order items: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id for order items: %w", err)
	}
	oi.ID = id
	return nil
}

func (s *MySQLStorer) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if fn returns an error

	if err := fn(tx); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
func (s *MySQLStorer) ListOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	err := s.db.SelectContext(ctx, &orders, "SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// Fetch order items for each order
	for i := range orders {
		var items []OrderItem
		err = s.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id = ?", orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items for order %d: %w", orders[i].ID, err)
		}
		orders[i].Items = items
	}
	return orders, nil
}

//UpdateOrder

// DeleteOrder
func (s *MySQLStorer) GetOrder(ctx context.Context, userId int64) (*Order, error) {
	var o Order
	err := s.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE user_id = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Fetch order items
	var items []OrderItem
	err = s.db.SelectContext(ctx, &o.Items, "SELECT * FROM order_items WHERE order_id = ?", o.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	o.Items = items
	return &o, nil
}

func (s *MySQLStorer) DeleteOrder(ctx context.Context, id int64) error {
	// Start a transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if any error occurs

	// Delete order items
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete order items: %w", err)
	}

	// Delete the order
	_, err = tx.ExecContext(ctx, "DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
func (s *MySQLStorer) CreateUser(ctx context.Context, u *User) (*User, error) {
	query := `INSERT INTO users (name, email, password, is_admin) VALUES (:name, :email, :password, :is_admin)`
	res, err := s.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	u.ID = id
	return u, nil
}
func (s *MySQLStorer) GetUser(ctx context.Context, email string) (*User, error) {
	var u User
	err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}
func (s *MySQLStorer) UpdateUser(ctx context.Context, u *User) (*User, error) {
	query := `
		UPDATE users SET
			name = :name,
			email = :email,
			password = :password,
			is_admin = :is_admin,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := s.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return u, nil
}
func (s *MySQLStorer) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
func (s *MySQLStorer) ListUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := s.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
func (s *MySQLStorer) CreateSession(ctx context.Context, session *Session) (*Session, error) {
	_, err := s.db.NamedExecContext(ctx, `INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES (:id, :user_email, :refresh_token, :is_revoked, :expires_at)`, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return session, nil
}
func (s *MySQLStorer) GetSession(ctx context.Context, id string) (*Session, error) {
	var session Session
	err := s.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}
func (s *MySQLStorer) RevokeSession(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE sessions SET is_revoked = TRUE WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}
	return nil
}
func (s *MySQLStorer) DeleteSession(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
