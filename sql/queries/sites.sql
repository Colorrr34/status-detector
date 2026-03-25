-- name: CreateSite :one
INSERT INTO sites (name,url,description,created_at,updated_at)
VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: GetSites :many
SELECT * FROM sites;