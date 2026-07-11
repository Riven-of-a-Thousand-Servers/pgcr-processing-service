-- name: CreatePlayerRaidStats :exec
INSERT INTO player_raid_stats (
    raid_name,
    raid_difficulty,
    player_membership_id,
    kills,
    deaths,
    assists,
    hour_played,
    clears,
    full_clears,
    flawless,
    contest_clear,
    day_one,
    solo,
    duo,
    trio,
    solo_flawless,
    duo_flawless,
    trio_flawless)
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
    $16,
    $17,
    $18)
ON CONFLICT (
    raid_name,
    raid_difficulty,
    player_membership_id)
    DO UPDATE SET
        kills = player_raid_stats.kills + EXCLUDED.kills,
        deaths = player_raid_stats.deaths + EXCLUDED.deaths,
        assists = player_raid_stats.assists + EXCLUDED.assists,
        hour_played = player_raid_stats.hour_played + EXCLUDED.hour_played,
        clears = player_raid_stats.clears + EXCLUDED.clears,
        full_clears = player_raid_stats.full_clears + EXCLUDED.full_clears,
        flawless = player_raid_stats.flawless
        OR EXCLUDED.flawless,
        contest_clear = player_raid_stats.contest_clear
        OR EXCLUDED.contest_clear,
        day_one = player_raid_stats.day_one
        OR EXCLUDED.day_one,
        solo = player_raid_stats.solo
        OR EXCLUDED.solo,
        duo = player_raid_stats.duo
        OR EXCLUDED.duo,
        trio = player_raid_stats.trio
        OR EXCLUDED.trio,
        solo_flawless = player_raid_stats.solo_flawless
        OR EXCLUDED.solo_flawless,
        duo_flawless = player_raid_stats.duo_flawless
        OR EXCLUDED.duo_flawless,
        trio_flawless = player_raid_stats.trio_flawless
        OR EXCLUDED.trio_flawless
