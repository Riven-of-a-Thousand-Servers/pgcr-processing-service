-- name: CreatePlayer :exec
INSERT INTO player (
    membership_id,
    membership_type,
    global_display_name,
    global_display_name_code,
    display_name,
    last_seen)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6)
ON CONFLICT (
    membership_id)
    DO UPDATE SET
        last_seen = EXCLUDED.last_seen
