-- name: Add :exec
INSERT INTO incidents (
        id, title, severity, service_name, started_at, resolved_at
        ) 
    VALUES (
        $1, $2, $3, $4, $5, $6
        );

-- name: Edit :exec
UPDATE incidents
    SET id = $1,
    title = $2,
    severity = $3,
    service_name = $4,
    started_at = $5,
    resolved_at = $6
WHERE id = $7
RETURNING *;

-- name: GetAll :many
SELECT * FROM incidents;

-- name: GetByID :one
SELECT * FROM incidents
WHERE id = $1;

-- name: DeleteByID :exec
DELETE FROM incidents
WHERE id = $1;

-- name: DeleteAll :exec
DELETE FROM incidents;