CREATE TYPE core.member_role AS ENUM ('default', 'admin');

CREATE TABLE IF NOT EXISTS core.group_chat_members (
	id SERIAL PRIMARY KEY,
	chat_uid UUID NOT NULL,
	user_uid UUID NOT NULL,
	role core.member_role NOT NULL DEFAULT 'default'
);