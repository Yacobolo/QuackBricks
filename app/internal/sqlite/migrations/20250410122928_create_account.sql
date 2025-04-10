-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS account (
    id TEXT PRIMARY KEY,                      -- UUID or logical user ID
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS account_group (
    id TEXT PRIMARY KEY,                      -- UUID or logical group ID
    name TEXT UNIQUE NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS account_group_membership (
    account_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    PRIMARY KEY (account_id, group_id),
    FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES account_group(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS account_group_membership;
DROP TABLE IF EXISTS account_group;
DROP TABLE IF EXISTS account;
-- +goose StatementEnd