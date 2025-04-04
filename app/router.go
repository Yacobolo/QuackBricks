package app

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
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

func setupRoutes(ctx context.Context, router chi.Router) (err error) {
	defer router.Handle("/static/*", hashfs.FileServer(staticSys))

	cfg := NewConfig()

	db, err := NewDuckDB(cfg)

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
	authHandler, err := NewAuthHandler(cfg.JWKSURL)

	if err != nil {
		log.Fatalf("Error creating AuthHandler: %s", err)
	}

	apiRouter.Use(AuthMiddleware(authHandler))

	if err := errors.Join(
		// setupHome(router, sessionSignals, ns, index),
		setupHome(router, apiRouter, db),
	); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	router.Mount("/api", apiRouter)

	return nil
}
