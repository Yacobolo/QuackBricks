-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS catalog (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,                        -- Logical name of the table
    schema_name TEXT NOT NULL,                        -- Logical schema
    source_type TEXT NOT NULL,                        -- e.g. 'parquet', 'csv', 'delta'
    location TEXT NOT NULL,                           -- Path to the file/folder
    description TEXT,                                 -- Optional docstring
    registered_at DATETIME NOT NULL  
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS catalog;
-- +goose StatementEnd
