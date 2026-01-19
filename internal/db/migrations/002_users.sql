-- +goose Up

create Table users(
  id UUID PRIMARY KEY NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL,
  hashed_password TEXT not null Default 'unset',
  email TEXT not null unique
);

-- +goose Down

drop Table users;