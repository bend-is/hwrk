CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE table IF NOT EXISTS events
(
    id          uuid primary key default uuid_generate_v4(),
    title       varchar(100) not null,
    description text         not null,
    user_id     bigint       not null,
    start_at    timestamp    not null,
    finish_at   timestamp    not null,
    notify_at   timestamp    not null
);

CREATE INDEX IF NOT EXISTS events_user_id_idx on events (user_id);
CREATE INDEX IF NOT EXISTS events_start_at_idx on events (start_at);
