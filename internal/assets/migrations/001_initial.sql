-- +migrate Up
create table blobs (
    id SERIAL PRIMARY KEY,
    owner_id TEXT,
    purpose TEXT,
    blob_content json not null
);
-- +migrate Down
drop table blobs;