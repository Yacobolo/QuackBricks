package catalog

import (
	"context"
	"duckdb-test/app/internal/sqlite"
	"duckdb-test/pkg/catalog"
	"slices"
)

type CatalogService struct {
	q *sqlite.Queries
}

func NewService(q *sqlite.Queries) *CatalogService {
	return &CatalogService{
		q: q,
	}
}

func (c *CatalogService) Register(ctx context.Context, req *catalog.CatalogEntryInput) error {
	if err := validate(req); err != nil {
		return err
	}

	params := sqlite.CreateCatalogEntryParams{
		Name:        req.Name,
		SourceType:  req.SourceType,
		Location:    req.Location,
		SchemaName:  req.SchemaName,
		Description: req.Description,
	}

	return c.q.CreateCatalogEntry(ctx, params)
}

func validate(req *catalog.CatalogEntryInput) error {
	if req.Name == "" || req.SourceType == "" || req.Location == "" {
		return catalog.ErrMissingField
	}

	if !slices.Contains(catalog.AllowedSourceTypes, catalog.SourceType(req.SourceType)) {
		return catalog.ErrInvalidSource
	}

	return nil
}
