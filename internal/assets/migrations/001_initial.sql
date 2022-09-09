-- +migrate Up
create table blobs (
    id SERIAL PRIMARY KEY,
    owner_address Bytea,
    purpose TEXT,
    blob_content json not null
);
create table documents(
    id SERIAL PRIMARY KEY,
    type TEXT,
    name TEXT,
    image_url TEXT,
    purpose TEXT,
    owner_address BYTEA
);
-- +migrate Down
drop table blobs;
drop table documents;