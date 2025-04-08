package web

import (
	"duckdb-test/app/internal/auth"
	"duckdb-test/app/internal/domain"
	"duckdb-test/app/internal/duckdb"
	"duckdb-test/app/internal/handler"
	"duckdb-test/app/internal/sqlite"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func setupHome(router chi.Router, apiRouter chi.Router, db duckdb.DuckDB, q *sqlite.Queries) error {

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

		err := handler.QueryHandler(q, db, query, w)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error processing query: %s", err), http.StatusInternalServerError)
			return
		}

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

	apiRouter.Route("/catalog", func(catalogRouter chi.Router) {
		catalogRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			catalog, err := q.ListCatalogEntries(r.Context())

			if err != nil {
				http.Error(w, fmt.Sprintf("Error listing catalog entries: %v", err), http.StatusInternalServerError)
				return
			}

			jsonData, err := json.Marshal(catalog)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)

		})

		catalogRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			if id == "" {
				http.Error(w, "Catalog ID is required", http.StatusBadRequest)
				return
			}

			catalog, err := q.GetCatalogEntry(r.Context(), id)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error getting catalog entry: %v", err), http.StatusInternalServerError)
				return
			}

			jsonData, err := json.Marshal(catalog)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)

		})
		catalogRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Name        string  `json:"name"`
				SourceType  string  `json:"source_type"`
				Location    string  `json:"location"`
				SchemaName  *string `json:"schema_name"`
				Description *string `json:"description"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			params := sqlite.CreateCatalogEntryParams{
				Name:        req.Name,
				SourceType:  req.SourceType,
				Location:    req.Location,
				SchemaName:  req.SchemaName,
				Description: req.Description,
			}

			if err := domain.ValidateCatalogParams(params); err != nil {
				http.Error(w, fmt.Sprintf("Invalid catalog entry: %v", err), http.StatusBadRequest)
				return
			}

			if err := q.CreateCatalogEntry(r.Context(), params); err != nil {
				http.Error(w, fmt.Sprintf("Error creating catalog entry: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Catalog entry created successfully"))
		})
	})

	return nil
}
