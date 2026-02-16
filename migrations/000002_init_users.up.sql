CREATE TABLE IF NOT EXISTS core.users (
	id SERIAL PRIMARY KEY,
	uid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
	username VARCHAR(30) UNIQUE,
	avatar_uid UUID
);