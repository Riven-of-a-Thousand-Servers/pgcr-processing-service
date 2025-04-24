CREATE TYPE RAID_NAME as ENUM (
    'Salvation''s Edge',
    'Crota''s End',
    'Root of Nightmares',
    'King''s Fall',
    'Vow of the Disciple',
    'Vault of Glass',
    'Deep Stone Crypt',
    'Garden of Salvation',
    'Crown of Sorrow',
    'Last Wish',
    'Leviathan, Spire of Stars',
    'Leviathan, Eater of Worlds',
    'Leviathan',
    'Scourge of the Past');

CREATE TYPE RAID_DIFFICULTY as ENUM (
    'Normal',
    'Master',
    'Prestige',
    'Guided Games',
    'Challenge Mode'
    );

CREATE TYPE CHARACTER_CLASS as ENUM (
    'Titan',
    'Warlock',
    'Hunter'
    );

CREATE TYPE DAMAGE_TYPE as ENUM (
    'Kinetic',
    'Arc',
    'Void',
    'Solar',
    'Stasis',
    'Strand'
    );

CREATE TYPE EQUIPMENT_SLOT as ENUM (
    'Primary',
    'Special',
    'Heavy'
    );

CREATE TABlE raid_pgcr
(
    instance_id BIGINT PRIMARY KEY,
    blob        BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS player
(
    membership_id            BIGINT PRIMARY KEY,
    membership_type          INTEGER,
    global_display_name      VARCHAR,
    global_display_name_code INTEGER,
    display_name             VARCHAR,
    last_seen                TIMESTAMP
);

CREATE INDEX player_display_name_id ON player (display_name);

CREATE TABLE IF NOT EXISTS player_character
(
    character_id         BIGINT PRIMARY KEY,
    character_class      CHARACTER_CLASS,
    current_emblem       BIGINT,
    player_membership_id BIGINT,
    CONSTRAINT membership_id_fk FOREIGN KEY (player_membership_id) REFERENCES player (membership_id)
);

CREATE TABLE IF NOT EXISTS raid
(
    raid_name       RAID_NAME,
    raid_difficulty RAID_DIFFICULTY,
    is_active       BOOLEAN,
    release_date    TIMESTAMP,
    CONSTRAINT raid_pk PRIMARY KEY (raid_name, raid_difficulty)
);

CREATE TABLE IF NOT EXISTS raid_hash
(
    raid_hash       BIGINT PRIMARY KEY,
    raid_name       RAID_NAME       NOT NULL,
    raid_difficulty RAID_DIFFICULTY NOT NULL,
    CONSTRAINT raid_fk FOREIGN KEY (raid_name, raid_difficulty) REFERENCES raid (raid_name, raid_difficulty) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS player_raid_stats
(
    raid_name            RAID_NAME,
    raid_difficulty      RAID_DIFFICULTY,
    player_membership_id BIGINT,
    kills                INTEGER DEFAULT 0,
    melee_kills          INTEGER DEFAULT 0,
    super_kills          INTEGER DEFAULT 0,
    grenade_kills        INTEGER DEFAULT 0,
    deaths               INTEGER DEFAULT 0,
    assists              INTEGER DEFAULT 0,
    hour_played          INTEGER DEFAULT 0,
    clears               INTEGER DEFAULT 0,
    full_clears          INTEGER DEFAULT 0,
    flawless             BOOLEAN DEFAULT false,
    contest_clear        BOOLEAN DEFAULT false,
    day_one              BOOLEAN DEFAULT false,
    solo                 BOOLEAN DEFAULT false,
    duo                  BOOLEAN DEFAULT false,
    trio                 BOOLEAN DEFAULT false,
    solo_flawless        BOOLEAN DEFAULT false,
    duo_flawless         BOOLEAN DEFAULT false,
    trio_flawless        BOOLEAN DEFAULT false,
    PRIMARY KEY (raid_name, raid_difficulty, player_membership_id),
    CONSTRAINT raid_name_fk FOREIGN KEY (raid_name, raid_difficulty) REFERENCES raid (raid_name, raid_difficulty),
    CONSTRAINT membership_id_fk FOREIGN KEY (player_membership_id) REFERENCES player (membership_id)
);

CREATE TABLE IF NOT EXISTS instance_activity_stats
(
    instance_id          BIGINT,
    player_membership_id BIGINT,
    player_character_id  BIGINT,
    character_emblem     BIGINT,
    is_completed         BOOLEAN,
    kills                INTEGER,
    deaths               INTEGER,
    assists              INTEGER,
    melee_kills          INTEGER,
    grenade_kills        INTEGER,
    super_kills          INTEGER,
    kills_deaths_assists FLOAT,
    kills_deaths_ratio   FLOAT,
    efficiency           FLOAT,
    duration_seconds     INTEGER,
    time_played_seconds  INTEGER,
    CONSTRAINT player_raid_activity_stats_pk PRIMARY KEY (instance_id, player_membership_id, player_character_id),
    CONSTRAINT instance_id_fk FOREIGN KEY (instance_id) REFERENCES raid_pgcr (instance_id),
    CONSTRAINT player_membership_id_fk FOREIGN KEY (player_membership_id) REFERENCES player (membership_id),
    CONSTRAINT player_character_id_fk FOREIGN KEY (player_character_id) REFERENCES player_character (character_id)
)

CREATE TABLE IF NOT EXISTS weapon
(
    weapon_hash           BIGINT PRIMARY KEY,
    weapon_icon           VARCHAR,
    weapon_name           VARCHAR,
    weapon_damage_type    DAMAGE_TYPE,
    weapon_equipment_slot EQUIPMENT_SLOT
);
CREATE UNIQUE INDEX weapon_name_idx ON weapon (weapon_name);
CREATE UNIQUE INDEX weapon_equipment_slot_idx ON weapon (weapon_equipment_slot);

CREATE TABLE IF NOT EXISTS instance_activity_weapon_stats
(
    instance_id           BIGINT,
    player_character_id   BIGINT,
    weapon_id             BIGINT,
    total_kills           INTEGER,
    total_precision_kills INTEGER,
    precision_ratio       FLOAT,
    CONSTRAINT instance_activity_weapon_stats_pk PRIMARY KEY (instance_id, player_character_id, weapon_id),
    CONSTRAINT instance_id_fk FOREIGN KEY (instance_id) REFERENCES raid_pgcr (instance_id),
    CONSTRAINT player_character_id_fk FOREIGN KEY (player_character_id) REFERENCES player_character (character_id),
    CONSTRAINT weapon_id_fk FOREIGN KEY (weapon_id) REFERENCES weapon (weapon_hash)
);
CREATE UNIQUE INDEX instance_activity_weapon_stats_total_kills_idx ON instance_activity_weapon_stats (total_kills);
