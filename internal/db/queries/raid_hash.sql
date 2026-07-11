-- name: CreateRaidHash :exec
INSERT INTO raid_hash (
    raid_hash,
    raid_name,
    raid_difficulty)
VALUES (
    $1,
    $2,
    $3)
ON CONFLICT
    DO NOTHING
