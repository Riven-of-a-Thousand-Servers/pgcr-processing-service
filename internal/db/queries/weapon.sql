-- name: CreateWeapon :exec
INSERT INTO weapon (
    weapon_hash,
    icon_url,
    weapon_name,
    damage_type,
    equipment_slot
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);
