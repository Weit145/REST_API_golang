package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type Order struct {
	ID    int
	Name  string
	Price float64
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot open sqlite db: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: cannot ping sqlite db: %w", op, err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("%s: cannot migrate sqlite db: %w", op, err)
	}

	return &Storage{db: db}, nil

}

func migrate(db *sql.DB) error {

	const op = "storage.sqlite.migrate"

	schema := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		price FLOAT NOT NULL DEFAULT 0
	);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("%s: cannot execute migration: %w", op, err)
	}

	return nil
}
