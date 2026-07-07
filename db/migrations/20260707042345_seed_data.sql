-- +goose Up
INSERT INTO activity_name(
	id,
	label,
	release_date,
	is_active
)
VALUES(
	1,
	'Leviathan',
	'2017-09-13 9:00:00 -7:00',
	FALSE
),
(
	2,
	'Leviathan, Eater of Worlds',
	'2017-12-08 09:00:00 -8:00',
	FALSE
),
(
	3,
	'Leviathan, Spire of Stars',
	'2018-05-11 09:00:00 -7:00',
	FALSE
),
(
	4,
	'Last Wish',
	'2018-09-18 9:00:00 -7:00',
	TRUE
),
(
	5,
	'Scourge of the Past',
	'2018-12-07 09:00:00 -8:00',
	FALSE
),
(
	6,
	'Crown of Sorrow',
	'2019-06-04 09:00:00 -7:00',
	FALSE
),
(
	7,
	'Garden of Salvation',
	'2019-10-05 09:00:00 -7:00',
	TRUE
),
(
	8,
	'Deep Stone Crypt',
	'2020-11-21 09:00:00 -8:00',
	TRUE
),
(
	9,
	'Vault of Glass',
	'2021-05-22 09:00:00 -7:00',
	TRUE
),
(
	10,
	'Vow of the Disciple',
	'2022-03-05 09:00:00 -7:00',
	TRUE
),
(
	11,
	'King's Fall',
	'2022-08-26 09:00:00 -8:00',
	TRUE
),
(
	12,
	'Root of Nightmares',
	'2023-03-10 09:00:00 -8:00',
	TRUE
),
(
	13,
	'Crota's End',
	'2023-09-01 09:00:00 -8:00',
	TRUE
),
(
	14,
	'The Pantheon',
	'2024-04-30 09:00:00 -8:00',
	FALSE
),
(
	15,
	'Salvation's Edge',
	'2024-06-07 09:00:00 -7:00',
	TRUE
);
INSERT INTO activity_difficulty(
	id,
	label
)
VALUES(
	1,
	'Normal'
),
(
	2,
	'Guided Games'
),
(
	3,
	'Prestige'
),
(
	4,
	'Master'
),
(
	5,
	'Contest'
),
(
	6,
	'Tempo's Edge'
),
(
	7,
	'Regicide'
),
(
	8,
	'Superior Swordplay'
),
(
	9,
	'Atraks Sovereign'
),
(
	10,
	'Oryx Exalted'
),
(
	11,
	'Rhulk Indomitable'
),
(
	12,
	'Nezarec Sublime'
);
INSERT INTO activity(
	hash,
	label,
	name_id,
	difficulty_id,
	is_worlds_first
)
VALUES(
	2693136601,
	'Leviathan',
	1,
	1,
	TRUE
),
(
	2693136600,
	'Leviathan',
	1,
	1,
	FALSE
),
(
	2693136602,
	'Leviathan',
	1,
	1,
	FALSE
),
(
	2693136603,
	'Leviathan',
	1,
	1,
	FALSE
),
(
	2693136604,
	'Leviathan',
	1,
	1,
	FALSE
),
(
	2693136605,
	'Leviathan',
	1,
	1,
	FALSE
),
(
	89727599,
	'Leviathan: Guided Games',
	1,
	2,
	FALSE
),
(
	287649202,
	'Leviathan: Guided Games',
	1,
	2,
	FALSE
),
(
	1699948563,
	'Leviathan: Guided Games',
	1,
	2,
	FALSE
),
(
	1875726950,
	'Leviathan: Guided Games',
	1,
	2,
	FALSE
),
(
	3916343513,
	'Leviathan: Guided Games',
	1,
	2,
	FALSE
),
(
	4039317196,
	'Leviathan: Guided Games',
	1,
	2,
	FALSE
),
(
	417231112,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	508802457,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	757116822,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	771164842,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	1685065161,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	1800508819,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	2449714930,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	3446541099,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	4206123728,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	3912437239,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	3879860661,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	3857338478,
	'Leviathan: Prestige',
	1,
	3,
	FALSE
),
(
	3089205900,
	'Leviathan, Eater of Worlds',
	2,
	1,
	TRUE
),
(
	2164432138,
	'Leviathan, Eater of Worlds: Guided Games',
	2,
	2,
	FALSE
),
(
	809170886,
	'Leviathan, Eater of Worlds: Prestige',
	2,
	3,
	FALSE
),
(
	119944200,
	'Leviathan, Spire of Stars',
	3,
	1,
	TRUE
),
(
	3004605630,
	'Leviathan, Spire of Stars: Guided Games',
	3,
	2,
	FALSE
),
(
	3213556450,
	'Leviathan, Spire of Stars: Prestige',
	3,
	3,
	FALSE
),
(
	2122313384,
	'Last Wish',
	4,
	1,
	TRUE
),
(
	2214608157,
	'Last Wish',
	4,
	1,
	FALSE
),
(
	1661734046,
	'Last Wish: Guided Games',
	4,
	2,
	FALSE
),
(
	548750096,
	'Scourge of the Past',
	5,
	1,
	TRUE
),
(
	2812525063,
	'Scourge of the Past: Guided Games',
	5,
	2,
	FALSE
),
(
	3333172150,
	'Crown of Sorrow',
	6,
	1,
	TRUE
),
(
	960175301,
	'Crown of Sorrow: Guided Games',
	6,
	2,
	FALSE
),
(
	2659723068,
	'Garden of Salvation',
	7,
	1,
	TRUE
),
(
	3458480158,
	'Garden of Salvation',
	7,
	1,
	FALSE
),
(
	1042180643,
	'Garden of Salvation',
	7,
	1,
	FALSE
),
(
	2497200493,
	'Garden of Salvation: Guided Games',
	7,
	2,
	FALSE
),
(
	3845997235,
	'Garden of Salvation: Guided Games',
	7,
	2,
	FALSE
),
(
	3823237780,
	'Garden of Salvation: Guided Games',
	7,
	2,
	FALSE
),
(
	910380154,
	'Deep Stone Crypt',
	8,
	1,
	TRUE
),
(
	3976949817,
	'Deep Stone Crypt: Guided Games',
	8,
	2,
	FALSE
),
(
	3881495763,
	'Vault of Glass',
	9,
	1,
	FALSE
),
(
	3711931140,
	'Vault of Glass: Guided Games',
	9,
	2,
	FALSE
),
(
	1485585878,
	'Vault of Glass: Temp's Edge',
	9,
	6,
	TRUE
),
(
	1681562271,
	'Vault of Glass: Master',
	9,
	4,
	FALSE
),
(
	3022541210,
	'Vault of Glass: Master',
	9,
	4,
	FALSE
),
(
	1441982566,
	'Vow of the Disciple',
	10,
	1,
	TRUE
),
(
	2906950631,
	'Vow of the Discple',
	10,
	1,
	FALSE
),
(
	4156879541,
	'Vow of the Disciple: Guided Games',
	10,
	2,
	FALSE
),
(
	4217492330,
	'Vow of the Disciple: Master',
	10,
	4,
	FALSE
),
(
	3889634515,
	'Vow of the Disciple: Master',
	10,
	4,
	TRUE
),
(
	1374392663,
	'King's Fall',
	11,
	1,
	FALSE
),
(
	2897223272,
	'King's Fall: Guided Games',
	11,
	2,
	FALSE
),
(
	1063970578,
	'King's Fall: Regicide',
	11,
	7,
	TRUE
),
(
	2964135793,
	'King's Fall: Master',
	11,
	4,
	FALSE
),
(
	3257594522,
	'King's Fall: Master',
	11,
	4,
	FALSE
),
(
	2381413764,
	'Root of Nightmares',
	12,
	1,
	TRUE
),
(
	1191701339,
	'Root of Nightmares: Guided Games',
	12,
	2,
	FALSE
),
(
	2918919505,
	'Root of Nightmares: Master',
	12,
	4,
	FALSE
),
(
	4179289725,
	'Crota's End',
	13,
	1,
	FALSE
),
(
	107319834,
	'Crota's End',
	13,
	1,
	FALSE
),
(
	4103176774,
	'Crota's End: Guided Games',
	13,
	2,
	FALSE
),
(
	156253568,
	'Crota's End: Superior Swordplay',
	13,
	8,
	TRUE
),
(
	1507509200,
	'Crota's End: Master',
	13,
	4,
	FALSE
),
(
	4169648179,
	'The Pantheon: Atraks Sovereign',
	14,
	9,
	NULL
),
(
	4169648176,
	'The Pantheon: Oryx Exalted',
	14,
	10,
	NULL
),
(
	4169648177,
	'The Pantheon: Rhulk Indomitable',
	14,
	11,
	NULL
),
(
	4169648182,
	'The Pantheon: Nezarec Sublime',
	14,
	12,
	NULL
),
(
	1541433876,
	'Salvation's Edge',
	15,
	1,
	FALSE
),
(
	2192826039,
	'Salvation's Edge: Contest',
	15,
	5,
	TRUE
),
(
	4129614942,
	'Salvation's Edge: Master',
	15,
	4,
	FALSE
);
