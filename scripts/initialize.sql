CREATE TABLE IF NOT EXISTS guilds (
	id BIGINT PRIMARY KEY NOT NULL,
	locale VARCHAR(20) NOT NULL,
	color INTEGER NOT NULL,
	flags BIGINT NOT NULL,
  max_warns INTEGER,
  final_warn_punishment INTEGER,
  final_warn_punishment_duration BIGINT
);

CREATE TABLE IF NOT EXISTS autoroles (
	guild_id BIGINT NOT NULL,
	role_id BIGINT NOT NULL
);
