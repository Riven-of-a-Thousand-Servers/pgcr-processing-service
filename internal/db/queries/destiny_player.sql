-- name: CreateDestinyPlayer :exec
INSERT INTO destiny_player (
    membership_id,
    membership_type,
    icon_path,
    display_name,
    global_display_name,
    global_display_name_code,
    total_clears,
    total_full_clears,
    is_private,
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
);
