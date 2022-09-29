-- +migrate Up
CREATE TABLE blobs (
    id BIGSERIAL PRIMARY KEY,
    owner_address BYTEA,
    purpose TEXT,
    blob_content JSON NOT NULL
);
CREATE TABLE documents(
    id BIGSERIAL PRIMARY KEY,
    type TEXT,
    name TEXT,
    purpose TEXT,
    owner_address BYTEA
);
-- +migrate Down
DROP TABLE blobs;
DROP TABLE documents;