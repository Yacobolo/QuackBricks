package duckdb

import (
	"duckdb-test/app/internal/config"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/marcboeker/go-duckdb/v2"
)

type DuckDB interface {
	Query(query string) (*sqlx.Rows, error)
}

type duckDB struct {
	db *sqlx.DB
}

func New(cfg *config.Config) (DuckDB, error) {
	db, err := sqlx.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	q := fmt.Sprintf(`CREATE SECRET secret1 (TYPE azure, CONNECTION_STRING '%s');`, cfg.ConnectionString)
	if _, err := db.Exec(q); err != nil {
		db.Close()
		return nil, fmt.Errorf("create secret: %w", err)
	}

	return &duckDB{db}, nil
}

func (d *duckDB) Query(query string) (*sqlx.Rows, error) {
	return d.db.Queryx(query)
}
