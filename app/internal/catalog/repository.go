package catalog

import "context"

type CatalogRepository interface {
	GetCatalogEntryByID(ctx context.Context, id string) (*CatalogEntry, error)
	GetCatalogEntryByQualifiedName(ctx context.Context, name, schema string) (*CatalogEntry, error)
	ListCatalogEntries(ctx context.Context) ([]*CatalogEntry, error)
	RegisterCatalogEntry(ctx context.Context, entry *CatalogEntry) error
	DeleteCatalogEntry(ctx context.Context, id string) error
	UpdateCatalogEntry(ctx context.Context, entry *CatalogEntry) error
}
