package web

import (
	auth "duckdb-test/app/internal/auth/app"
	"duckdb-test/app/internal/catalog"
	"duckdb-test/app/internal/duckdb"
	"duckdb-test/app/internal/handler"
	"duckdb-test/app/internal/sqlite"
	"duckdb-test/app/internal/sqlite/gen"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func setupHome(router chi.Router, apiRouter chi.Router, db duckdb.DuckDB, q *gen.Queries) error {

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

		catalogRepo := sqlite.NewCatalogRepository(q)
		catalogService := catalog.NewService(catalogRepo)
		catalogHandler := catalog.NewHandler(catalogService)

		catalogRouter.Get("/", catalogHandler.ListCatalogEntries)
		catalogRouter.Post("/", catalogHandler.RegisterCatalogEntry)
		catalogRouter.Get("/{id}", catalogHandler.GetCatalogEntry)
		catalogRouter.Delete("/{id}", catalogHandler.DeleteCatalogEntry)
	})

	return nil
}
