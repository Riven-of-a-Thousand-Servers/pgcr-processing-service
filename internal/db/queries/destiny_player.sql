-- name: CreateDestinyPlayer :one
INSERT INTO destiny_player (
    membership_id,
    membership_type,
    icon_path,
    display_name,
    global_display_name,
    global_display_name_code,
    total_clears,
    total_full_clears,
    is_public,
    last_crawled,
    last_seen
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11
) ON CONFLICT (membership_id)
DO UPDATE
    SET
        membership_id = excluded.membership_id,
        membership_type = excluded.membership_type,
        global_display_name = excluded.global_display_name,
        global_display_name_code = excluded.global_display_name_code,
        icon_path = excluded.icon_path,
        is_public = excluded.is_public,
        last_seen = now(),
        last_crawled = now()
RETURNING *;

-- name: IncrementPlayerCounts :exec
UPDATE destiny_player
SET
    total_clears = total_clears + CASE WHEN $2::boolean THEN 1 ELSE 0 END,
    total_full_clears = total_full_clears + CASE WHEN $3::boolean THEN 1 ELSE 0 END
WHERE membership_id = $1;
