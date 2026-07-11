-- name: CreateInstanceActivity :exec
INSERT INTO instance_activity_stats (
    instance_id,
    player_membership_id,
    player_character_id,
    character_emblem,
    is_completed,
    kills,
    deaths,
    assists,
    kills_deaths_assists,
    kills_deaths_ratio,
    duration_seconds,
    time_played_seconds)
VALUES (
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
    $11,
    $12,
    $13)
