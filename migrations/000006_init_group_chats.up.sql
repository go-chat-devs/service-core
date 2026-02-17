CREATE TABLE IF NOT EXISTS core.group_chats (
	uid					UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title				VARCHAR(128) NOT NULL,
	bio					VARCHAR(512),
	avatar_uid	UUID,
	created_at	TIMESTAMPTZ NOT NULL DEFAULT NOW()
);