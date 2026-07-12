-- name: CreateInstancePlayer :exec
INSERT INTO instance_player (
    instance_id,
    membership_id,
    completed,
    time_played_seconds
) VALUES (
    $1,
    $2,
    $3,
    $4
) ON CONFLICT (instance_id, membership_id) DO NOTHING;
