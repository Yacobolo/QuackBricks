package catalog

import (
	"errors"

	"github.com/google/uuid"
)

var ErrTableNotFound = errors.New("table not found")

type Schema struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tables      []*Table  `json:"tables"`
}

type CreateSchemaParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateSchemaParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func newSchema(params *CreateSchemaParams) (*Schema, error) {
	schema := &Schema{
		ID:          uuid.New(),
		Name:        params.Name,
		Description: params.Description,
	}

	if err := schema.validate(); err != nil {
		return nil, err
	}

	return schema, nil
}

func (s *Schema) update(params *UpdateSchemaParams) error {
	if params.Name != "" {
		s.Name = params.Name
	}

	if params.Description != "" {
		s.Description = params.Description
	}

	if err := s.validate(); err != nil {
		return err
	}

	return nil
}

func (s *Schema) validate() error {
	if s.Name == "" {
		return ErrMissingField
	}

	return nil
}

func (s *Schema) createTable(catalogID uuid.UUID, params *CreateTableParams) (*Table, error) {
	table, err := newTable(catalogID, s.ID, params)
	if err != nil {
		return nil, err
	}

	s.Tables = append(s.Tables, table)

	return table, nil
}

func (s *Schema) deleteTable(tableID string) error {
	idx, _, err := s.findTableByID(tableID)
	if err != nil {
		return ErrTableNotFound
	}

	s.Tables = append(s.Tables[:*idx], s.Tables[*idx+1:]...)

	return nil
}

func (s *Schema) updateTable(tableID string, params *UpdateTableParams) (*Table, error) {
	_, table, err := s.findTableByID(tableID)
	if err != nil {
		return nil, ErrTableNotFound
	}

	if err := table.update(params); err != nil {
		return nil, err
	}

	return table, nil
}

// Helper
func (s *Schema) findTableByID(tableID string) (*int, *Table, error) {
	for idx, table := range s.Tables {
		if table.ID.String() == tableID {
			return &idx, table, nil
		}
	}

	return nil, nil, ErrTableNotFound
}
