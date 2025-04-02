package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/marcboeker/go-duckdb"
)

func main() {

	cfg := NewConfig()

	authHandler, err := NewAuthHandler(cfg.JWKSURL)

	if err != nil {
		log.Fatalf("Error creating AuthHandler: %s", err)
	}

	db, err := NewDuckDB(cfg)

	if err != nil {
		log.Fatalf("Error creating DuckDB: %s", err)
	}

	// Create a logger
	logger := log.New(os.Stdout, "[app] ", log.LstdFlags)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(apiRouter chi.Router) {

		apiRouter.Use(AuthMiddleware(authHandler))

		apiRouter.Get("/query", func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("q")

			if query == "" {
				http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
				return
			}

			logger.Println("Received query:", query)

			queryHandler(db, query, w, logger)

		})

		apiRouter.Get("/me", func(w http.ResponseWriter, r *http.Request) {

			user := GetUserFromContext(r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			jsonData, err := json.Marshal(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)

		})
	})

	logger.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func queryHandler(db DuckDB, query string, w http.ResponseWriter, logger *log.Logger) {
	rows, err := db.Query(query)
	if err != nil {
		logger.Printf("Error executing query: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		rowData := make(map[string]interface{})
		err = rows.MapScan(rowData)
		if err != nil {
			logger.Printf("Error scanning row: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, rowData)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(results)
	if err != nil {
		logger.Printf("Error marshaling JSON: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
	logger.Printf("Query executed successfully: %s", query)
}
