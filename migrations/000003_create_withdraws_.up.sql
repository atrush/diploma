CREATE TABLE IF NOT EXISTS withdraws
(
    id           uuid primary key default uuid_generate_v4(),
    user_id      uuid not null,
    "number"     varchar,
    uploaded_at  TIMESTAMP,
    "sum"      int,
    FOREIGN KEY (user_id) REFERENCES users (id)
);