-- +migrate Up

CREATE TABLE blobs (
    id BIGSERIAL PRIMARY KEY,
    owner_address TEXT,
    purpose TEXT,
    blob_content JSON NOT NULL
);

CREATE TABLE documents(
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    file_key TEXT,
    owner_address TEXT,
    mime_type TEXT
);

-- +migrate Down

DROP TABLE blobs;
DROP TABLE documents;