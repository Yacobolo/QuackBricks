package catalog

import (
	"slices"

	"github.com/google/uuid"
)

type Table struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Schema      *Schema    `json:"schema"`
	Catalog     *Catalog   `json:"catalog"`
	SourceType  SourceType `json:"source_type"`
	Location    string     `json:"location"`
}

type CreateTableParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Schema      *Schema  `json:"schema"`
	Catalog     *Catalog `json:"catalog"`
}

func NewTable(params *CreateTableParams) (*Table, error) {
	t := &Table{
		ID:          uuid.New(),
		Name:        params.Name,
		Description: params.Description,
		Schema:      params.Schema,
		Catalog:     params.Catalog,
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (c *Table) Validate() error {
	if c.Name == "" || c.SourceType == "" || c.Location == "" {
		return ErrMissingField
	}

	if !slices.Contains(AllowedSourceTypes, SourceType(c.SourceType)) {
		return ErrInvalidSource
	}

	if c.Schema == nil {
		return ErrMissingField
	}

	if c.Catalog == nil {
		return ErrMissingField
	}

	return nil
}
