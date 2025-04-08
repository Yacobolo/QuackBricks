package domain

import (
	"duckdb-test/app/internal/sqlite"
	"errors"
)

var (
	ErrMissingField    = errors.New("missing required field")
	ErrInvalidSource   = errors.New("invalid source_type")
	allowedSourceTypes = map[string]bool{
		"parquet": true,
		"delta":   true,
	}
)

func ValidateCatalogParams(params sqlite.CreateCatalogEntryParams) error {
	if params.Name == "" || params.SourceType == "" || params.Location == "" {
		return ErrMissingField
	}
	if !allowedSourceTypes[params.SourceType] {
		return ErrInvalidSource
	}

	return nil
}
