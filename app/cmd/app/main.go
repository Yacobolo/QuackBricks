package main

import (
	"context"
	"duckdb-test/app"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/delaneyj/toolbelt"
)

const (
	port = 8080
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting Server", "url", fmt.Sprintf("http://localhost:%d", port))
	defer logger.Info("Stopping Server")

	ctx := context.Background()

	if err := run(ctx); err != nil {
		logger.Error("Error running server", slog.Any("err", err))
		os.Exit(1)
	}

}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg := toolbelt.NewErrGroupSharedCtx(
		ctx,
		app.RunBlocking(port, nil),
	)
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}

	return nil
}
