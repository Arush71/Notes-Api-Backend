-- name: CreateUser :one

insert Into users(id,updated_at, hashed_password , email)
values (
  gen_random_uuid(),
  NOW(),
  $2,
  $1
)
returning
id,
email,
created_at;

-- name: FindUserByEmail :one

Select id, hashed_password, created_at
from users where email = $1;