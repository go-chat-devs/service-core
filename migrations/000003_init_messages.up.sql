CREATE TYPE core.message_type AS ENUM ('text', 'image');

CREATE TABLE IF NOT EXISTS core.messages (
	id SERIAL PRIMARY KEY,
	chat_uid UUID NOT NULL,
	type core.message_type NOT NULL,
	sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS core.text_messages (
	uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	message_id INT REFERENCES core.messages(id) ON DELETE CASCADE,
	content TEXT,
	user_uid UUID,
	changed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS core.image_messages (
	uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	message_id INT REFERENCES core.messages(id) ON DELETE CASCADE,
	file_uid UUID NOT NULL DEFAULT gen_random_uuid(),
	user_uid UUID
);