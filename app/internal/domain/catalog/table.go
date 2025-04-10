package catalog

import (
	"errors"
	"slices"

	"github.com/google/uuid"
)

var ErrMissingCatalog = errors.New("missing catalog")
var ErrMissingSchema = errors.New("missing schema")

type SourceType string

const (
	SourceTypeParquet SourceType = "parquet"
	SourceTypeDelta   SourceType = "delta"
)

var AllowedSourceTypes = []SourceType{
	SourceTypeParquet,
	SourceTypeDelta,
}

type Table struct {
	ID          uuid.UUID  `json:"id"`
	CatalogID   uuid.UUID  `json:"catalog_id"`
	SchemaID    uuid.UUID  `json:"schema_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SourceType  SourceType `json:"source_type"`
	Location    string     `json:"location"`
}

type CreateTableParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Schema      *Schema  `json:"schema"`
	Catalog     *Catalog `json:"catalog"`
}

type UpdateTableParams struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SourceType  SourceType `json:"source_type"`
	Location    string     `json:"location"`
}

func newTable(catalogID uuid.UUID, schemaID uuid.UUID, params *CreateTableParams) (*Table, error) {
	t := &Table{
		ID:          uuid.New(),
		Name:        params.Name,
		Description: params.Description,
		CatalogID:   catalogID,
		SchemaID:    schemaID,
	}

	if err := t.validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Table) update(params *UpdateTableParams) error {
	if params.Name != "" {
		t.Name = params.Name
	}
	if params.Description != "" {
		t.Description = params.Description
	}
	if params.SourceType != "" {
		t.SourceType = params.SourceType
	}
	if params.Location != "" {
		t.Location = params.Location
	}

	if err := t.validate(); err != nil {
		return err
	}

	return nil
}

func (c *Table) validate() error {
	if c.Name == "" || c.SourceType == "" || c.Location == "" {
		return ErrMissingField
	}

	if !slices.Contains(AllowedSourceTypes, SourceType(c.SourceType)) {
		return ErrInvalidSource
	}

	if c.CatalogID == uuid.Nil {
		return ErrMissingCatalog
	}

	if c.SchemaID == uuid.Nil {
		return ErrMissingSchema
	}

	return nil
}
