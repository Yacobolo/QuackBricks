package catalog

import "github.com/google/uuid"

type Schema struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Catalog     *Catalog  `json:"catalog"`
}

type CreateSchemaParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Catalog     *Catalog `json:"catalog"`
}

func NewSchema(params *CreateSchemaParams) (*Schema, error) {
	schema := &Schema{
		ID:          uuid.New(),
		Name:        params.Name,
		Description: params.Description,
		Catalog:     params.Catalog,
	}

	if err := schema.Validate(); err != nil {
		return nil, err
	}

	return schema, nil
}

func (s *Schema) Validate() error {
	if s.Name == "" {
		return ErrMissingField
	}

	if s.Catalog == nil {
		return ErrMissingField
	}

	return nil
}
