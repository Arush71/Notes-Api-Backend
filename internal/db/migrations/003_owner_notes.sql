-- +goose Up

alter table notes
add column owner_id UUID references users(id) ON DELETE cascade;

-- +goose Down

alter table users
drop column owner_id;