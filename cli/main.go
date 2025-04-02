/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"duckdb-test/cli/cmd"
	"duckdb-test/cli/internal/config"
)

func main() {

	cfg := config.NewConfig()

	cmd.Execute(cfg)
}
