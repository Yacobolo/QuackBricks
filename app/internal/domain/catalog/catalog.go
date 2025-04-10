package catalog

import (
	"github.com/google/uuid"
)

type Catalog struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
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
