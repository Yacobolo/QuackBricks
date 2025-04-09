package main

import (
	"duckdb-test/app/internal/cli"
	"duckdb-test/app/internal/config"
)

func main() {

	cfg := config.NewConfig()

	cli.Execute(cfg)
}
