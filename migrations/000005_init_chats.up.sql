CREATE TABLE IF NOT EXISTS core.chats (
	uid						UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_uid_low	UUID,
	user_uid_high	UUID,
	created_at		TIMESTAMPTZ NOT NULL DEFAULT NOW()
);