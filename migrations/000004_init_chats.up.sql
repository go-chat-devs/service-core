CREATE TABLE IF NOT EXISTS chats (
	id SERIAL PRIMARY KEY,
	uid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
	user_uid_low UUID,
	user_uid_high UUID
);