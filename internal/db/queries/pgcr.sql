-- name: CreatePgcr :exec
INSERT INTO pgcr (
    instance_id,
    blob
)
VALUES (
    $1,
    $2
);
