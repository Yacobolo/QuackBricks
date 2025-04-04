package web

import (
	"duckdb-test/app/internal/auth"
	"duckdb-test/app/internal/duckdb"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func setupHome(router chi.Router, apiRouter chi.Router, db duckdb.DuckDB) error {

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// sse := datastar.NewSSE(w, r)

		c := home()
		c.Render(r.Context(), w)

		// sse.MergeFragmentTempl(c)
	})

	apiRouter.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		if query == "" {
			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		fmt.Println("Received query:", query)

		queryHandler(db, query, w)

	})

	apiRouter.Get("/me", func(w http.ResponseWriter, r *http.Request) {

		user := auth.GetUserFromContext(r)

		jsonData, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)

	})

	return nil
}

func queryHandler(db duckdb.DuckDB, query string, w http.ResponseWriter) {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error executing query: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		rowData := make(map[string]interface{})
		err = rows.MapScan(rowData)
		if err != nil {
			fmt.Printf("Error scanning row: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, rowData)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(results)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
	fmt.Printf("Query executed successfully: %s", query)
}
