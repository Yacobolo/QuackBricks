package catalog

import (
	"errors"
)

var (
	ErrInvalidSource = errors.New("invalid source_type")
	ErrMissingField  = errors.New("missing required field")
	ErrInvalidInput  = errors.New("invalid input")
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
