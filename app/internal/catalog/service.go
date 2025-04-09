package catalog

import (
	"context"
	"duckdb-test/app/internal/sqlite"
	"duckdb-test/pkg/catalog"
	"slices"
)

type Service struct {
	q *sqlite.Queries
}

func NewService(q *sqlite.Queries) *Service {
	return &Service{
		q: q,
	}
}

func (c *Service) Register(ctx context.Context, req *catalog.CatalogEntryInput) error {
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

func (c *Service) ListCatalogEntries(ctx context.Context) ([]sqlite.Catalog, error) {
	catalogEntries, err := c.q.ListCatalogEntries(ctx)
	if err != nil {
		return nil, err
	}

	return catalogEntries, nil
}

func (c *Service) DeleteCatalogEntry(ctx context.Context, name string) error {
	return c.q.DeleteCatalogEntry(ctx, name)
}

func (c *Service) GetCatalogEntry(ctx context.Context, name string) (sqlite.Catalog, error) {
	catalogEntry, err := c.q.GetCatalogEntry(ctx, name)
	if err != nil {
		return sqlite.Catalog{}, err
	}

	return catalogEntry, nil
}
