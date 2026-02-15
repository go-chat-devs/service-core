CREATE TYPE MEMBER_ROLE AS ENUM ('default', 'admin');

CREATE TABLE IF NOT EXISTS group_chat_members (
	id SERIAL PRIMARY KEY,
	chat_uid UUID NOT NULL,
	user_uid UUID NOT NULL,
	role MEMBER_ROLE NOT NULL DEFAULT 'default'
);