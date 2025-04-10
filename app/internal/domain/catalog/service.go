package catalog

import (
	"context"
	"fmt"
)

type Service struct {
	repo CatalogRepository
}

func NewService(repo CatalogRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// create catalog

func (c *Service) CreateCatalog(ctx context.Context, params *CreateCatalogParams) (*Catalog, error) {
	catalog, err := NewCatalog(params)
	if err != nil {
		return nil, fmt.Errorf("error creating catalog: %w", err)
	}
	err = c.repo.CreateCatalog(ctx, catalog)
	if err != nil {
		return nil, fmt.Errorf("error creating catalog: %w", err)
	}
	return catalog, nil
}

func (c *Service) GetCatalog(ctx context.Context, id string) (*Catalog, error) {
	catalog, err := c.repo.GetCatalogById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting catalog: %w", err)
	}
	return catalog, nil
}

func (c *Service) DeleteCatalog(ctx context.Context, id string) error {
	err := c.repo.DeleteCatalog(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting catalog: %w", err)
	}
	return nil
}

func (c *Service) ListCatalogs(ctx context.Context) ([]*Catalog, error) {
	catalogs, err := c.repo.ListCatalogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing catalogs: %w", err)
	}
	return catalogs, nil
}

// Schema

func (c *Service) CreateSchema(ctx context.Context, catalogID string, params *CreateSchemaParams) (*Schema, error) {
	catalog, err := c.repo.GetCatalogById(ctx, catalogID)
	if err != nil {
		return nil, fmt.Errorf("error getting catalog: %w", err)
	}

	schema, err := catalog.CreateSchema(params)
	if err != nil {
		return nil, fmt.Errorf("error creating schema: %w", err)
	}
	err = c.repo.CreateSchema(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("error creating schema: %w", err)
	}
	return schema, nil
}

func (c *Service) UpdateSchema(ctx context.Context, catalogID string, schemaID string, params *UpdateSchemaParams) (*Schema, error) {
	catalog, err := c.repo.GetCatalogById(ctx, catalogID)
	if err != nil {
		return nil, fmt.Errorf("error getting catalog: %w", err)
	}

	schema, err := catalog.UpdateSchema(schemaID, params)
	if err != nil {
		return nil, fmt.Errorf("error updating schema: %w", err)
	}

	err = c.repo.UpdateSchema(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("error updating schema: %w", err)
	}
	return schema, nil
}

func (c *Service) GetSchema(ctx context.Context, id string) (*Schema, error) {
	schema, err := c.repo.GetSchemaByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting schema: %w", err)
	}
	return schema, nil
}

func (c *Service) DeleteSchema(ctx context.Context, id string) error {
	err := c.repo.DeleteSchema(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting schema: %w", err)
	}
	return nil
}
func (c *Service) ListSchemas(ctx context.Context) ([]*Schema, error) {
	schemas, err := c.repo.ListSchemas(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing schemas: %w", err)
	}
	return schemas, nil
}

// func (c *Service) Register(ctx context.Context, req *CreateCatalogEntryParams) error {
// 	catalogEntry, err := NewCatalogEntry(req)

// 	if err != nil {
// 		return fmt.Errorf("error creating catalog entry: %w", err)
// 	}

// 	err = c.repo.RegisterCatalogEntry(ctx, catalogEntry)
// 	if err != nil {
// 		return fmt.Errorf("error registering catalog entry: %w", err)
// 	}

// 	return nil
// }

// func (c *Service) ListCatalogEntries(ctx context.Context) ([]*CatalogEntry, error) {
// 	catalogEntries, err := c.repo.ListCatalogEntries(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("error listing catalog entries: %w", err)
// 	}
// 	return catalogEntries, nil
// }

// func (c *Service) DeleteCatalogEntry(ctx context.Context, id string) error {
// 	err := c.repo.DeleteCatalogEntry(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("error deleting catalog entry: %w", err)
// 	}

// 	return nil
// }

// func (c *Service) GetCatalogEntry(ctx context.Context, nameOrId string) (*CatalogEntry, error) {
// 	nameOrId = strings.TrimSpace(nameOrId)

// 	// Check for qualified name (e.g., schema.name)
// 	if strings.Contains(nameOrId, ".") {
// 		schemaName, name, ok := strings.Cut(nameOrId, ".")
// 		if !ok || schemaName == "" || name == "" {
// 			return nil, fmt.Errorf("invalid qualified name: %s", nameOrId)
// 		}

// 		entry, err := c.repo.GetCatalogEntryByQualifiedName(ctx, name, schemaName)
// 		if err != nil {
// 			return nil, fmt.Errorf("error getting catalog entry by qualified name: %w", err)
// 		}
// 		return entry, nil
// 	}

// 	// Check for valid UUID
// 	if _, err := uuid.Parse(nameOrId); err != nil {
// 		return nil, fmt.Errorf("invalid UUID %s: %w", nameOrId, err)
// 	}

// 	entry, err := c.repo.GetCatalogEntryByID(ctx, nameOrId)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting catalog entry by ID: %w", err)
// 	}

// 	return entry, nil
// }

// //update

// func (c *Service) UpdateCatalogEntry(ctx context.Context, id string, params *CreateCatalogEntryParams) error {
// 	catalogEntry, err := c.repo.GetCatalogEntryByID(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("error getting catalog entry: %w", err)
// 	}

// 	// Update the fields of the catalog entry
// 	if params.Name != "" {
// 		catalogEntry.Name = params.Name
// 	}
// 	if params.SourceType != "" {
// 		catalogEntry.SourceType = params.SourceType
// 	}
// 	if params.Location != "" {
// 		catalogEntry.Location = params.Location
// 	}
// 	if params.SchemaName != "" {
// 		catalogEntry.SchemaName = params.SchemaName
// 	}
// 	if params.Description != nil {
// 		catalogEntry.Description = params.Description
// 	}

// 	if err := catalogEntry.Validate(); err != nil {
// 		return fmt.Errorf("error validating catalog entry: %w", err)
// 	}

// 	err = c.repo.UpdateCatalogEntry(ctx, catalogEntry)
// 	if err != nil {
// 		return fmt.Errorf("error updating catalog entry: %w", err)
// 	}

// 	return nil
// }
