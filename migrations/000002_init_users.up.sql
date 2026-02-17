CREATE TABLE IF NOT EXISTS core.users (
	uid					UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	username		VARCHAR(30) UNIQUE,
	avatar_uid	UUID
);