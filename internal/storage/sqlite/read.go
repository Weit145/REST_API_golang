package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Weit145/REST_API_golang/internal/storage"
)

func (s *Storage) ReadOrder(name string) (Order, error) {
	const op = "storage.sqlite.readOrder"
	var order Order

	stmt, err := s.db.Prepare("SELECT id, name, price FROM orders WHERE name = ?")
	if err != nil {
		return order, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(name).Scan(&order.ID, &order.Name, &order.Price)
	if errors.Is(err, sql.ErrNoRows) {
		return order, fmt.Errorf("%s: order with name %s not found: %w", op, name, storage.ErrURLNotFound)
	}
	if err != nil {
		return order, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}
