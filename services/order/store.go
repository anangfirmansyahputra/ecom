package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anangfirmansyahp5/ecom/services/auth"
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

func (s *Store) CreateOrder(ctx context.Context, payload types.CreateOrderPayload) error {
	userID, ok := ctx.Value(auth.UserKey).(int)

	if !ok || userID <= 0 {
		return errors.New("invalid or missing userID")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction :%v", err)
	}

	res, err := tx.Exec(`
		INSERT INTO orders (userId, total, status, address, createdAt)
		VALUES (?, ?, ?, ?, NOW())
	`, userID, 0, "pending", "")

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert order : %v", err)
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get order ID: %v", err)
	}

	totalAmount := 0.0
	for _, item := range payload.Items {
		var price float64
		var quantity int64
		err := s.db.QueryRow(`
			SELECT price, quantity FROM products where id = ?`, item.ProductID).Scan(&price, &quantity)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get product price: %v", err)
		}

		itemTotal := price * float64(item.Quantity)
		totalAmount += itemTotal

		if item.Quantity > int(quantity) {
			tx.Rollback()
			return fmt.Errorf("stock of this product is less than the quantity of your order")
		}

		_, err = tx.Exec(`
			INSERT INTO order_items (orderId, productId, quantity, price)
			VALUES (?, ?, ?, ?)`,
			orderID, item.ProductID, item.Quantity, price,
		)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert order item: %v", err)
		}

		_, err = tx.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", item.Quantity, item.ProductID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update product quantity: %v", err)
		}

	}

	_, err = tx.Exec(`
		UPDATE orders SET total = ? WHERE id = ?`,
		totalAmount, orderID,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update order total: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
