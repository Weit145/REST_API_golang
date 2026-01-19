package sqlite

import (
	"fmt"

	"github.com/Weit145/REST_API_golang/internal/storage"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) CreateOrder(order Order) error {
	const op = "storage.sqlite.createOrder"

	stmt, err := s.db.Prepare(
		"INSERT INTO orders (name, price) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(order.Name, order.Price)
	if err != nil {
		if sqlite3Err, ok := err.(sqlite3.Error); ok && sqlite3Err.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: order with name %s already exists: %w", op, order.Name, storage.ErrURLExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
