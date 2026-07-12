-- name: CreateInstanceCharacter :exec
INSERT INTO instance_character (
    instance_id,
    membership_id,
    character_id,
    emblem_hash,
    class_hash,
    completed,
    kills,
    deaths,
    assists,
    kda,
    kdr,
    super_kills,
    melee_kills,
    grenade_kills,
    efficiency,
    time_played_seconds
)
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
    $13,
    $14,
    $15,
    $16
)
ON CONFLICT (instance_id, membership_id, character_id) DO NOTHING
RETURNING instance_id;
