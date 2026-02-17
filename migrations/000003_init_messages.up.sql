CREATE TYPE core.message_type AS ENUM ('text', 'image');

CREATE TABLE IF NOT EXISTS core.messages (
	uid					UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	chat_uid		UUID NOT NULL,
	sender_uid	UUID,
	type 				core.message_type NOT NULL,
	sent_at 		TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS core.text_messages (
	id					SERIAL PRIMARY KEY,
	message_uid	UUID NOT NULL REFERENCES core.messages(uid) ON DELETE CASCADE,
	content			TEXT NOT NULL,
	changed_at	TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS core.image_messages (
	id					SERIAL PRIMARY KEY,
	message_uid	UUID NOT NULL REFERENCES core.messages(uid) ON DELETE CASCADE,
	file_uid		UUID NOT NULL DEFAULT gen_random_uuid(),
);