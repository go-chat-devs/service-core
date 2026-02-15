CREATE TYPE message_type AS ENUM ('text', 'image');

CREATE TABLE IF NOT EXISTS messages(
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL UNIQUE,
    chat_uid UUID NOT NULL,
    type message_type NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS text_messages(
    id SERIAL PRIMARY KEY,
    message_uid UUID NOT NULL UNIQUE,
    text TEXT,
    from UUID,
    changed TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS image_messages(
    id SERIAL PRIMARY KEY,
    message_uid UUID NOT NULL UNIQUE,
    file_uid UUID NOT NULL DEFAULT gen_random_uuid(),
    from UUID
);