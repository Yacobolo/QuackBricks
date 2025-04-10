package catalog

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type CatalogEntry struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	SchemaName   string    `json:"schema_name"`
	SourceType   string    `json:"source_type"`
	Location     string    `json:"location"`
	Description  *string   `json:"description"`
	RegisteredAt time.Time `json:"registered_at"`
}

type CreateCatalogEntryParams struct {
	Name        string  `json:"name"`
	SourceType  string  `json:"source_type"`
	Location    string  `json:"location"`
	SchemaName  string  `json:"schema_name"`
	Description *string `json:"description"`
}

// new
func NewCatalogEntry(params *CreateCatalogEntryParams) (*CatalogEntry, error) {
	c := &CatalogEntry{
		ID:           uuid.New(),
		Name:         params.Name,
		SchemaName:   params.SchemaName,
		SourceType:   params.SourceType,
		Location:     params.Location,
		Description:  params.Description,
		RegisteredAt: time.Now(),
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CatalogEntry) Validate() error {
	if c.Name == "" || c.SourceType == "" || c.Location == "" {
		return ErrMissingField
	}

	if !slices.Contains(AllowedSourceTypes, SourceType(c.SourceType)) {
		return ErrInvalidSource
	}

	return nil
}
