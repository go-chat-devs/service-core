CREATE TYPE  message_type AS ENUM ('text', 'image');


CREATE TABLE IF NOT EXISTS messages(
    id SERIAL PRIMARY KEY,
    uid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    chat_uid UUID NOT NULL DEFAULT gen_random_uuid(),
    timestamp TIMESTAMP NOT NULL,
    sender_uid UUID NOT NULL DEFAULT gen_random_uuid(),
    type message_type NOT NULL
);


CREATE TABLE IF NOT EXISTS text_messages(
    id SERIAL PRIMARY KEY,
    message_uid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    text TEXT,
    changed TIMESTAMP
);


CREATE TABLE IF NOT EXISTS image_messages(
    id SERIAL PRIMARY KEY,
    message_uid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    file_uid UUID NOT NULL DEFAULT gen_random_uuid()
);