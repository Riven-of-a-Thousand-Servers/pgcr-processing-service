-- name: CreateInstanceCharacterWeapon :exec
INSERT INTO instance_character_weapon (
    instance_id,
    player_membership_id,
    player_character_id,
    weapon_id,
    kills,
    precision_kills,
    precision_ratio
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT DO NOTHING
RETURNING instance_id;
