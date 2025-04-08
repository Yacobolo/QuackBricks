-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS catalog (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,                        -- Logical name of the table
    source_type TEXT NOT NULL,                        -- e.g. 'parquet', 'csv', 'delta'
    location TEXT NOT NULL,                           -- Path to the file/folder
    schema_name TEXT DEFAULT 'main',                  -- Logical schema (optional)
    description TEXT,                                 -- Optional docstring
    registered_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- Auto timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS catalog;
-- +goose StatementEnd
