package handler

import (
	"context"
	"duckdb-test/app/internal/duckdb"
	"duckdb-test/app/internal/sqlite"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strings"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/auxten/postgresql-parser/pkg/walk"
)

// Import the tree package for AST node types

func QueryHandler(q *sqlite.Queries, db duckdb.DuckDB, query string, w http.ResponseWriter) error {

	catalogMap := getCatalogMap(q)

	// Call the rewriteSQLTableSources function
	sql, err := rewriteSQLTableSources(query, catalogMap, nil)
	if err != nil {
		return fmt.Errorf("failed to rewrite SQL: %w", err)
	}

	json, err := db.QueryToJSON(sql)
	if err != nil {
		return fmt.Errorf("failed to query to JSON: %w", err)
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	_, err = w.Write(json)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil

}

func getCatalogMap(q *sqlite.Queries) map[string]string {
	ctx := context.Background()

	catalog, err := q.ListCatalogEntries(ctx)

	if err != nil {
		log.Printf("Error listing catalog entries: %v", err)
		return nil
	}

	catalogMap := make(map[string]string)
	for _, entry := range catalog {
		catalogMap[entry.Name] = fmt.Sprintf("delta_scan('%s')", entry.Location)
	}
	return catalogMap
}

func rewriteSQLTableSources(sql string, catalog map[string]string, logger *log.Logger) (string, error) {

	hashedCatalog := make(map[string]string)

	w := &walk.AstWalker{
		Fn: func(ctx interface{}, node interface{}) (stop bool) {

			// log.Printf("node type %T", node)
			// log.Printf("node value %s", node)

			switch node := node.(type) {
			case *tree.TableName:
				fmt.Printf(">>> Found TableName Node! %s\n", node)
				replaceTableNameWithHash(node, catalog, hashedCatalog)
			}
			return false

		},
	}

	stmts, err := parser.Parse(sql)
	if err != nil {
		return "", fmt.Errorf("failed to parse SQL: %w", err)
	}

	_, _ = w.Walk(stmts, nil)

	sql_out := stmts.String()

	for key, value := range hashedCatalog {
		sql_out = strings.ReplaceAll(sql_out, key, value)
	}

	return sql_out, nil

}

func replaceTableNameWithHash(expr *tree.TableName, catalog, hashedCatalog map[string]string) {
	nameStr := expr.TableName.String()

	value, ok := catalog[nameStr]
	if ok {
		hash := hashString(value)
		hashedCatalog[hash] = value
		expr.TableName = tree.Name(hash)
	}
	return
}

func hashString(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum64())
}
