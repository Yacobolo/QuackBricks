package web

import (
	"context"
	"duckdb-test/app/internal/auth"
	"duckdb-test/app/internal/config"
	"duckdb-test/app/internal/duckdb"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

func setupRoutes(ctx context.Context, router chi.Router) (err error) {
	defer router.Handle("/static/*", hashfs.FileServer(staticSys))

	cfg := config.NewConfig()

	db, err := duckdb.New(cfg)

	if err != nil {
		log.Fatalf("Error creating DuckDB: %s", err)
	}

	natsPort, err := toolbelt.FreePort()
	if err != nil {
		return fmt.Errorf("error getting free port: %w", err)
	}

	ns, err := embeddednats.New(ctx, embeddednats.WithNATSServerOptions(&natsserver.Options{
		JetStream: true,
		Port:      natsPort,
		StoreDir:  "data/nats",
	}))
	if err != nil {
		return fmt.Errorf("error creating embedded nats server: %w", err)
	}
	ns.WaitForServer()

	sessionSignals := sessions.NewCookieStore([]byte("datastar-session-secret"))
	sessionSignals.MaxAge(int(24 * time.Hour / time.Second))

	// auth
	apiRouter := chi.NewRouter()
	authHandler, err := auth.NewAuthHandler(cfg.JWKSURL)

	if err != nil {
		log.Fatalf("Error creating AuthHandler: %s", err)
	}

	apiRouter.Use(auth.AuthMiddleware(authHandler))

	if err := errors.Join(
		// setupHome(router, sessionSignals, ns, index),
		setupHome(router, apiRouter, db),
	); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	router.Mount("/api", apiRouter)

	return nil
}
