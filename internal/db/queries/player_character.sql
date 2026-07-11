-- name: CreatePlayerCharacter :exec
INSERT INTO player_character (
    character_id,
    character_class,
    current_emblem,
    player_membership_id)
VALUES (
    $1,
    $2,
    $3,
    $4)
ON CONFLICT (
    character_id)
    DO UPDATE SET
        character_emblem = EXCLUDED.current_emblem
