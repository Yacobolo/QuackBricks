package catalog

import (
	"errors"

	"github.com/google/uuid"
)

var ErrSchemaNotFound = errors.New("schema not found")

type Catalog struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Schemas     []*Schema `json:"schemas"`
}

type CreateCatalogParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// new
func NewCatalog(params *CreateCatalogParams) (*Catalog, error) {
	c := &Catalog{
		ID:          uuid.New(),
		Name:        params.Name,
		Description: params.Description,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Catalog) Validate() error {

	if c.Name == "" {
		return ErrMissingField
	}

	return nil
}

// SCHEMA

func (c *Catalog) CreateSchema(params *CreateSchemaParams) (*Schema, error) {
	schema, err := newSchema(params)
	if err != nil {
		return nil, err
	}

	c.Schemas = append(c.Schemas, schema)

	return schema, nil
}

func (c *Catalog) DeleteSchema(schemaID string) error {

	idx, _, err := c.findSchemaByID(schemaID)
	if err != nil {
		return ErrSchemaNotFound
	}

	// remove schema from catalog
	c.Schemas = append(c.Schemas[:*idx], c.Schemas[*idx+1:]...)

	return nil
}

func (c *Catalog) UpdateSchema(schemaID string, params *UpdateSchemaParams) (*Schema, error) {
	_, schema, err := c.findSchemaByID(schemaID)
	if err != nil {
		return nil, err
	}

	if err := schema.update(params); err != nil {
		return nil, err
	}

	return schema, nil
}

func (c *Catalog) findSchemaByID(schemaID string) (*int, *Schema, error) {
	for idx, schema := range c.Schemas {
		if schema.ID.String() == schemaID {
			return &idx, schema, nil
		}
	}
	return nil, nil, ErrSchemaNotFound
}

//table

func (c *Catalog) CreateTable(schemaID string, params *CreateTableParams) (*Table, error) {
	_, schema, err := c.findSchemaByID(schemaID)
	if err != nil {
		return nil, err
	}

	table, err := schema.createTable(c.ID, params)
	if err != nil {
		return nil, err
	}

	return table, nil
}

func (c *Catalog) DeleteTable(schemaID string, tableID string) error {
	schemaIdx, _, err := c.findSchemaByID(schemaID)
	if err != nil {
		return err
	}

	err = c.Schemas[*schemaIdx].deleteTable(tableID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Catalog) UpdateTable(schemaID string, tableID string, params *UpdateTableParams) (*Table, error) {
	schemaIdx, _, err := c.findSchemaByID(schemaID)
	if err != nil {
		return nil, err
	}

	table, err := c.Schemas[*schemaIdx].updateTable(tableID, params)
	if err != nil {
		return nil, err
	}

	return table, nil
}
