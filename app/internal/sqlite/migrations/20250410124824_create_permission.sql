-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS permission (
    id TEXT PRIMARY KEY,
    account_id TEXT REFERENCES account(id),
    group_id TEXT REFERENCES account_group(id),

    catalog_id TEXT REFERENCES catalog(id),
    schema_id TEXT REFERENCES schema(id),
    table_id TEXT REFERENCES table_registry(id),

    action TEXT NOT NULL, -- e.g., 'SELECT', 'MODIFY', 'ALL'
    created_at DATETIME NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permission;
-- +goose StatementEnd