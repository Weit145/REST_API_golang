package sqlite

import (
	"fmt"

	"github.com/Weit145/REST_API_golang/internal/storage"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) UpdateOrder(newOrder Order) error {
	const op = "storage.sqlite.updateOrder"

	stmt, err := s.db.Prepare("UPDATE orders SET name = ?, price = ? WHERE name = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(newOrder.Name, newOrder.Price, newOrder.Name)
	if err != nil {
		if sqlite3Err, ok := err.(sqlite3.Error); ok && sqlite3Err.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: order with name %s already exists: %w", op, newOrder.Name, storage.ErrURLExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
