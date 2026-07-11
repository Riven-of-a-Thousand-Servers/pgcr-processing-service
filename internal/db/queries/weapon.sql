-- name: CreateWeapon :exec
INSERT INTO weapon (
    weapon_hash,
    weapon_icon,
    weapon_name,
    weapon_damage_type,
    weapon_equipment_slot)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5)
