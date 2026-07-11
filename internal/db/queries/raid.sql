-- name: CreateRaid   :exec
INSERT INTO raid (
    raid_name,
    raid_difficulty,
    is_active,
    release_date)
VALUES (
    $1,
    $2,
    $3,
    $4)
ON CONFLICT
    DO NOTHING;

