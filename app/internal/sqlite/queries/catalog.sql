-- name: CreateCatalogEntry :exec
INSERT INTO catalog (
    id,
    name,
    schema_name,
    source_type,
    location,
    description,
    registered_at
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetCatalogEntryById :one
SELECT * FROM catalog
WHERE id = ?
LIMIT 1;

-- name: GetCatalogEntryByQualifiedName :one
SELECT * FROM catalog
WHERE name = ? and schema_name = ?
LIMIT 1;

-- name: ListCatalogEntries :many
SELECT * FROM catalog
ORDER BY registered_at DESC;

-- name: UpdateCatalogEntry :exec
UPDATE catalog
SET name = ?,
    schema_name = ?,
    source_type = ?,
    location = ?,
    description = ?,
    registered_at = ?
WHERE id = ?;


-- name: DeleteCatalogEntry :exec
DELETE FROM catalog
WHERE id = ?;