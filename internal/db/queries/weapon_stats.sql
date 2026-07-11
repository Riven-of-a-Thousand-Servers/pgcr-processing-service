-- name: CreateInstanceActivityWeaponStats :exec
INSERT INTO instance_activity_weapon_stats (
    instance_id,
    player_character_id,
    weapon_id,
    total_kills,
    total_precision_kills,
    precision_rates)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6)
ON CONFLICT
    DO NOTHING
