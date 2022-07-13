CREATE TABLE IF NOT EXISTS orders
(
    id           uuid primary key default uuid_generate_v4(),
    user_id      uuid not null,
    "number"     varchar,
    uploaded_at  TIMESTAMP,
    status       varchar,
    accrual      int,
    FOREIGN KEY (user_id) REFERENCES users (id)
);