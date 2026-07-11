-- +goose Up
CREATE TABLE IF NOT EXISTS pgcr(
	instance_id BIGINT PRIMARY KEY,
	BLOB BYTEA NOT NULL
);
-- +goose Down
DROP TABLE pgcr;
-- +goose Up
CREATE TABLE IF NOT EXISTS activity_name(
	id BIGINT PRIMARY KEY,
	label TEXT NOT NULL,
	is_active BOOLEAN NOT NULL,
	release_dateTIMESTAMP(0) WITH TIME ZONE NOT NULL
);
-- +goose Down
DROP TABLE activity_name;
-- +goose Up
CREATE TABLE IF NOT EXISTS activity_difficulty(
	id BIGINT PRIMARY KEY,
	label TEXT NOT NULL
);
-- +goose Down
DROP TABLE activity_difficulty;
-- +goose Up
CREATE TABLE IF NOT EXISTS activity(
	hash BIGINT PRIMARY KEY,
	label TEXT NOT NULL,
	name_id BIGINT NOT NULL,
	difficulty_id BIGINT NOT NULL,
	is_worlds_first BOOLEAN,
	CONSTRAINT activity_name_fk FOREIGN KEY(name_id) REFERENCES activity_name(id),
	CONSTRAINT activity_difficulty_fk FOREIGN KEY(difficulty_id) REFERENCES activity_difficulty(id)
);
-- +goose Down
DROP TABLE activity;
-- +goose Up
CREATE TABLE IF NOT EXISTS instance(
	id BIGINT PRIMARY KEY,
	activity_hash BIGINT NOT NULL,
	is_fresh BOOLEAN NOT NULL,
	flawless BOOLEAN NOT NULL,
	completed BOOLEAN NOT NULL,
	player_count INT NOT NULL,
	duration INT NOT NULL,
	end_time TIMESTAMP NOT NULL,
	start_Time TIMESTAMP NOT NULL CONSTRAINT instance_hash_fk FOREIGN KEY(activity_hash) REFERENCES activity(hash)
);
-- +goose Down
DROP TABLE instance;
-- +goose Up
CREATE TABLE IF NOT EXISTS destiny_player(
	membership_id BIGINT PRIMARY KEY,
	membership_type INT NOT NULL,
	icon_path TEXT,
	display_name TEXT,
	global_display_name TEXT,
	global_display_name_code INT,
	total_clears INT NOT NULL DEFAULT 0,
	total_full_clears INT NOT NULL DEFAULT 0,
	is_private BOOLEAN DEFAULT FALSE,
	last_crawled TIMESTAMP NOT NULL,
	last_seen TIMESTAMP
);
CREATE INDEX destiny_player_display_name_idx
ON destiny_player(display_name);
-- +goose Down
DROP TABLE destiny_player;
-- +goose Up
CREATE TABLE IF NOT EXISTS instance_player(instance_id BIGINT NOT NULL,
membership_id BIGINT NOT NULL,
completed BOOLEAN DEFAULT FALSE,
time_played_seconds INT NOT NULL CONSTRAINT instance_player_pk PRIMARY KEY(
	instance_id,
	membership_id
));
-- +goose Down
DROP TABLE instance_player;
-- +goose Up
CREATE TABLE IF NOT EXISTS instance_character(
	instance_id BIGINT NOT NULL,
	membership_id BIGINT NOT NULL,
	character_id BIGINT NOT NULL,
	class_hash TEXT NOT NULL,
	emblem_hash TEXT NOT NULL,
	completed BOOLEAN NOT NULL,
	kills INT NOT NULL,
	deaths INT NOT NULL,
	assists INT NOT NULL,
	kda DECIMAL NOT NULL,
	kdr DECIMAL NOT NULL,
	super_kills INT NOT NULL,
	melee_kills INT NOT NULL,
	grenade_kills INT NOT NULL,
	efficiency INT NOT NULL,
	time_played_seconds INT NOT NULL CONSTRAINT instance_character_pk PRIMARY KEY(instance_id,
	membership_id,
	character_id),
	CONSTRAINT instance_character_instance_fk FOREIGN KEY(instance_id) REFERENCES instance(id),
	CONSTRAINT instance_character_destiny_player_fk FOREIGN KEY(membership_id) REFERENCES destiny_player(membership_id)
);
-- +goose Down
DROP TABLE instance_character;
-- +goose Up
CREATE TABLE IF NOT EXISTS weapon(
	hash BIGINT PRIMARY KEY,
	icon_url TEXT NOT NULL,
	name TEXT NOT NULL,
	equipment_slot TEXT NOT NULL,
	damage_type TEXT NOT NULL
);
-- +goose Down
DROP TABLE weapon;
-- +goose Up
CREATE TABLE IF NOT EXISTS instance_character_weapon(
	instance_id BIGINT NOT NULL,
	player_membership_id BIGINT NOT NULL,
	player_character_id BIGINT NOT NULL,
	weapon_id BIGINT NOT NULL,
	kills INT NOT NULL DEFAULT 0,
	precision_kills INT NOT NULL DEFAULT 0,
	precision_ratio DECIMAL NOT NULL DEFAULT 0.0 CONSTRAINT instance_character_weapon_pk PRIMARY KEY(
		instance_id,
		player_membership_id,
		player_character_id,
		weapon_id
	),
	CONSTRAINT instance_character_weapon_instance_fk FOREIGN KEY(instance_id) REFERENCES instance(id),
	CONSTRAINT instance_character_weapon_player_fk FOREIGN KEY(player_membership_id) REFERENCES destiny_player(membership_id),
	CONSTRAINT instance_character_weapon_character_fk FOREIGN KEY(
		instance_id,
		player_membership_id,
		player_character_id
	) REFERENCES instance_character(
		instance_id,
		membership_id,
		character_id
	),
	CONSTRAINT instance_character_weapon_weapon_fk FOREIGN KEY(weapon_id) REFERENCES weapon(hash)
);
-- +goose Down
DROP TABLE instance_character_weapon;
