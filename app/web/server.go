package web

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/benbjohnson/hashfs"
	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var staticFS embed.FS

var (
	staticSys = hashfs.NewFS(staticFS)
)

func RunBlocking(port int, readyCh chan struct{}) toolbelt.CtxErrFunc {
	return func(ctx context.Context) error {

		router := chi.NewRouter()

		router.Use(
			middleware.Recoverer,
			middleware.Logger,
		)

		if err := setupRoutes(ctx, router); err != nil {
			return fmt.Errorf("error setting up routes: %w", err)
		}

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}

		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		if readyCh != nil {
			close(readyCh)
		}
		return srv.ListenAndServe()
	}
}
