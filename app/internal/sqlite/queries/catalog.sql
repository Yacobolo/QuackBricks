-- name: CreateCatalogEntry :exec
INSERT INTO catalog (
    name, source_type, location, schema_name, description
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetCatalogEntry :one
SELECT * FROM catalog
WHERE name = ?
LIMIT 1;

-- name: ListCatalogEntries :many
SELECT * FROM catalog
ORDER BY registered_at DESC;

-- name: UpdateCatalogEntry :exec
UPDATE catalog
SET
    source_type = ?,
    location = ?,
    schema_name = ?,
    description = ?
WHERE name = ?;

-- name: DeleteCatalogEntry :exec
DELETE FROM catalog
WHERE name = ?;