package sqlite

import (
	"context"
	"duckdb-test/app/internal/catalog"
	"duckdb-test/app/internal/sqlite/gen"

	"github.com/google/uuid"
)

type CatalogRepository struct {
	q *gen.Queries
}

func NewCatalogRepository(q *gen.Queries) catalog.CatalogRepository {
	return &CatalogRepository{
		q: q,
	}
}

func (r *CatalogRepository) GetCatalogEntryByID(ctx context.Context, id string) (*catalog.CatalogEntry, error) {
	entry, err := r.q.GetCatalogEntryById(ctx, id)
	if err != nil {
		return nil, err
	}

	return catalogToDomain(entry)

}

func (r *CatalogRepository) GetCatalogEntryByQualifiedName(ctx context.Context, name, schema string) (*catalog.CatalogEntry, error) {
	params := gen.GetCatalogEntryByQualifiedNameParams{
		Name:       name,
		SchemaName: schema,
	}
	entry, err := r.q.GetCatalogEntryByQualifiedName(ctx, params)
	if err != nil {
		return nil, err
	}

	return catalogToDomain(entry)
}

func (r *CatalogRepository) ListCatalogEntries(ctx context.Context) ([]*catalog.CatalogEntry, error) {
	entries, err := r.q.ListCatalogEntries(ctx)
	if err != nil {
		return nil, err
	}

	var catalogEntries []*catalog.CatalogEntry
	for _, entry := range entries {
		catalogEntry, err := catalogToDomain(entry)
		if err != nil {
			return nil, err
		}
		catalogEntries = append(catalogEntries, catalogEntry)
	}

	return catalogEntries, nil
}

func (r *CatalogRepository) RegisterCatalogEntry(ctx context.Context, entry *catalog.CatalogEntry) error {
	dbEntry := catalogToDB(entry)

	params := gen.CreateCatalogEntryParams{
		ID:           dbEntry.ID,
		Name:         dbEntry.Name,
		SourceType:   dbEntry.SourceType,
		Location:     dbEntry.Location,
		SchemaName:   dbEntry.SchemaName,
		Description:  dbEntry.Description,
		RegisteredAt: dbEntry.RegisteredAt,
	}

	err := r.q.CreateCatalogEntry(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *CatalogRepository) DeleteCatalogEntry(ctx context.Context, id string) error {
	err := r.q.DeleteCatalogEntry(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CatalogRepository) UpdateCatalogEntry(ctx context.Context, entry *catalog.CatalogEntry) error {
	dbEntry := catalogToDB(entry)

	params := gen.UpdateCatalogEntryParams{
		ID:           dbEntry.ID,
		Name:         dbEntry.Name,
		SourceType:   dbEntry.SourceType,
		Location:     dbEntry.Location,
		SchemaName:   dbEntry.SchemaName,
		Description:  dbEntry.Description,
		RegisteredAt: dbEntry.RegisteredAt,
	}

	err := r.q.UpdateCatalogEntry(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// helpers

func catalogToDomain(entry gen.Catalog) (*catalog.CatalogEntry, error) {

	id, err := uuid.Parse(entry.ID)
	if err != nil {
		return nil, err
	}

	return &catalog.CatalogEntry{
		ID:           id,
		Name:         entry.Name,
		SchemaName:   entry.SchemaName,
		SourceType:   entry.SourceType,
		Location:     entry.Location,
		Description:  entry.Description,
		RegisteredAt: entry.RegisteredAt,
	}, nil
}

func catalogToDB(entry *catalog.CatalogEntry) gen.Catalog {
	return gen.Catalog{
		ID:           entry.ID.String(),
		Name:         entry.Name,
		SchemaName:   entry.SchemaName,
		SourceType:   entry.SourceType,
		Location:     entry.Location,
		Description:  entry.Description,
		RegisteredAt: entry.RegisteredAt,
	}
}
