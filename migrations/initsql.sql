-- +goose Up
-- +goose StatementBegin
create table if not exists sessions
(
    rent_id       bigint  not null
                  constraint sessions_pk
                  primary key,
    user_id       bigint,
    started_at    timestamp,
    completed_at  timestamp,
    images_before text,
    images_after  text
);

create unique index if not exists sessions_rent_id_uindex
    on sessions (rent_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists sessions;
DROP index if exists sessions_rent_id_uindex;
-- +goose StatementEnd

CREATE USER 'user' WITH PASSWORD 'postgres' WITH SUPERUSER;