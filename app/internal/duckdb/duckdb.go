package duckdb

import (
	"duckdb-test/app/internal/config"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/marcboeker/go-duckdb/v2"
)

type DuckDB interface {
	Query(query string) (*sqlx.Rows, error)
	QueryToJSON(query string) ([]byte, error)
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

func (d *duckDB) QueryToJSON(query string) ([]byte, error) {
	rows, err := d.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		rowData := make(map[string]interface{})
		err = rows.MapScan(rowData)
		if err != nil {
			return nil, fmt.Errorf("map scan: %w", err)
		}
		results = append(results, rowData)
	}

	return json.Marshal(results)
}
