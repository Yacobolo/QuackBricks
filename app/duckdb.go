package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DuckDB interface {
	Query(query string) (*sqlx.Rows, error)

	// Meta Queries
	ListTables() (*sqlx.Rows, error)
	DescribeTable(table string) (*sqlx.Rows, error)
}

type duckDB struct {
	db *sqlx.DB
}

func NewDuckDB(cfg *config) (DuckDB, error) {
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

func (d *duckDB) ListTables() (*sqlx.Rows, error) {
	return d.db.Queryx("SHOW TABLES")
}

func (d *duckDB) DescribeTable(table string) (*sqlx.Rows, error) {
	return d.db.Queryx(fmt.Sprintf("DESCRIBE %s", table))
}
