package catalog

import "context"

type CatalogRepository interface {
	//catalog
	GetCatalogById(ctx context.Context, id string) (*Catalog, error)
	GetCatalogByName(ctx context.Context, name string) (*Catalog, error)
	DeleteCatalog(ctx context.Context, id string) error
	CreateCatalog(ctx context.Context, catalog *Catalog) error
	ListCatalogs(ctx context.Context) ([]*Catalog, error)
	UpdateCatalog(ctx context.Context, catalog *Catalog) error

	//schemas
	GetSchemaByID(ctx context.Context, id string) (*Schema, error)
	GetSchemaByName(ctx context.Context, name, catalog string) (*Schema, error)
	DeleteSchema(ctx context.Context, id string) error
	CreateSchema(ctx context.Context, schema *Schema) error
	ListSchemas(ctx context.Context) ([]*Schema, error)
	UpdateSchema(ctx context.Context, schema *Schema) error

	//table Registry
	GetTableByID(ctx context.Context, id string) (*Table, error)
	GetTableByQualifiedName(ctx context.Context, name, schema string) (*Table, error)
	ListTables(ctx context.Context) ([]*Table, error)
	RegisterTable(ctx context.Context, entry *Table) error
	DeleteTable(ctx context.Context, id string) error
	UpdateTable(ctx context.Context, entry *Table) error
}
