CREATE TYPE message_type AS ENUM ('text', 'image');


CREATE TABLE IF NOT EXISTS messages(
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    chat_uid VARCHAR(36) NOT NULL,
    timestamp INTEGER NOT NULL,
    sender_uid VARCHAR(36) NOT NULL,
    type message_type NOT NULL
);


CREATE TABLE IF NOT EXISTS text_messages(
    id SERIAL PRIMARY KEY,
    message_uid VARCHAR(36) NOT NULL UNIQUE,
    text VARCHAR(500)
);


CREATE TABLE IF NOT EXISTS image_messages(
    id SERIAL PRIMARY KEY,
    message_uid VARCHAR(36) NOT NULL UNIQUE,
    file_uid VARCHAR(36) NOT NULL
);