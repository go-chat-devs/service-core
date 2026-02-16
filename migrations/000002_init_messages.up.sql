CREATE TYPE MESSAGE_TYPE AS ENUM ('text', 'image');

CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	chat_uid UUID NOT NULL,
	type MESSAGE_TYPE NOT NULL,
	sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS text_messages (
	uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	message_id INT REFERENCES messages(id) ON DELETE CASCADE,
	content TEXT,
	user_uid UUID,
	changed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS image_messages (
	uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	message_id INT REFERENCES messages(id) ON DELETE CASCADE,
	file_uid UUID NOT NULL DEFAULT gen_random_uuid(),
	user_uid UUID
);