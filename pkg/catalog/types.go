package catalog

import "errors"

var (
	ErrInvalidSource = errors.New("invalid source_type")
	ErrMissingField  = errors.New("missing required field")
)

type SourceType string

const (
	SourceTypeParquet SourceType = "parquet"
	SourceTypeDelta   SourceType = "delta"
)

var AllowedSourceTypes = []SourceType{
	SourceTypeParquet,
	SourceTypeDelta,
}

type CatalogEntryInput struct {
	Name        string  `json:"name"`
	SourceType  string  `json:"source_type"`
	Location    string  `json:"location"`
	SchemaName  *string `json:"schema_name"`
	Description *string `json:"description"`
}
