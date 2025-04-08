// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlite

import (
	"time"
)

type Catalog struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	SourceType   string     `json:"source_type"`
	Location     string     `json:"location"`
	SchemaName   *string    `json:"schema_name"`
	Description  *string    `json:"description"`
	RegisteredAt *time.Time `json:"registered_at"`
}
