-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS catalog (
    id TEXT PRIMARY KEY,                        -- e.g., 'datamart'
    description TEXT
);

CREATE TABLE IF NOT EXISTS schema (
    id TEXT PRIMARY KEY,                        -- e.g., 'finance'
    catalog_id TEXT NOT NULL,
    description TEXT,
    FOREIGN KEY (catalog_id) REFERENCES catalog(id) ON DELETE CASCADE,
    UNIQUE (catalog_id, id)
);

CREATE TABLE IF NOT EXISTS table_registry (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,                        -- Logical name of the table
    schema_id TEXT NOT NULL,                        -- Logical schema
    source_type TEXT NOT NULL,                        -- e.g. 'parquet', 'csv', 'delta'
    location TEXT NOT NULL,                           -- Path to the file/folder
    description TEXT,                                 -- Optional docstring
    registered_at DATETIME NOT NULL  
    FOREIGN KEY (schema_id) REFERENCES schema(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS catalog;
DROP TABLE IF EXISTS schema;
DROP TABLE IF EXISTS table_registry;
-- +goose StatementEnd
