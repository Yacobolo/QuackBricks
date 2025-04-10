package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo CatalogRepository
}

func NewService(repo CatalogRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (c *Service) Register(ctx context.Context, req *CreateCatalogEntryParams) error {
	catalogEntry, err := NewCatalogEntry(req)

	if err != nil {
		return fmt.Errorf("error creating catalog entry: %w", err)
	}

	err = c.repo.RegisterCatalogEntry(ctx, catalogEntry)
	if err != nil {
		return fmt.Errorf("error registering catalog entry: %w", err)
	}

	return nil
}

func (c *Service) ListCatalogEntries(ctx context.Context) ([]*CatalogEntry, error) {
	catalogEntries, err := c.repo.ListCatalogEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing catalog entries: %w", err)
	}
	return catalogEntries, nil
}

func (c *Service) DeleteCatalogEntry(ctx context.Context, id string) error {
	err := c.repo.DeleteCatalogEntry(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting catalog entry: %w", err)
	}

	return nil
}

func (c *Service) GetCatalogEntry(ctx context.Context, nameOrId string) (*CatalogEntry, error) {
	nameOrId = strings.TrimSpace(nameOrId)

	// Check for qualified name (e.g., schema.name)
	if strings.Contains(nameOrId, ".") {
		schemaName, name, ok := strings.Cut(nameOrId, ".")
		if !ok || schemaName == "" || name == "" {
			return nil, fmt.Errorf("invalid qualified name: %s", nameOrId)
		}

		entry, err := c.repo.GetCatalogEntryByQualifiedName(ctx, name, schemaName)
		if err != nil {
			return nil, fmt.Errorf("error getting catalog entry by qualified name: %w", err)
		}
		return entry, nil
	}

	// Check for valid UUID
	if _, err := uuid.Parse(nameOrId); err != nil {
		return nil, fmt.Errorf("invalid UUID %s: %w", nameOrId, err)
	}

	entry, err := c.repo.GetCatalogEntryByID(ctx, nameOrId)
	if err != nil {
		return nil, fmt.Errorf("error getting catalog entry by ID: %w", err)
	}

	return entry, nil
}

//update

func (c *Service) UpdateCatalogEntry(ctx context.Context, id string, params *CreateCatalogEntryParams) error {
	catalogEntry, err := c.repo.GetCatalogEntryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting catalog entry: %w", err)
	}

	// Update the fields of the catalog entry
	if params.Name != "" {
		catalogEntry.Name = params.Name
	}
	if params.SourceType != "" {
		catalogEntry.SourceType = params.SourceType
	}
	if params.Location != "" {
		catalogEntry.Location = params.Location
	}
	if params.SchemaName != "" {
		catalogEntry.SchemaName = params.SchemaName
	}
	if params.Description != nil {
		catalogEntry.Description = params.Description
	}

	if err := catalogEntry.Validate(); err != nil {
		return fmt.Errorf("error validating catalog entry: %w", err)
	}

	err = c.repo.UpdateCatalogEntry(ctx, catalogEntry)
	if err != nil {
		return fmt.Errorf("error updating catalog entry: %w", err)
	}

	return nil
}
