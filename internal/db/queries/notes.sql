-- name: GetAllNotes :many

Select 
     id ,
    title,
    content,
    updated_at,
    owner_id
From notes where owner_id = $1;

-- name: CreateNewNote :one

Insert into notes (title, content , updated_at , owner_id)
Values ($1,$2 , NOW() , $3)
Returning owner_id , id , created_at;

-- name: GetANote :one

Select 
    id,
    title,
    content,
    updated_at
    owner_id
  from notes
where id = $1 and owner_id = $2;

-- name: DeleteNote :one
Delete from notes
where id = $1 and owner_id = $2 
returning id;

-- name: UpdateNote :one

update notes
Set 
    title = $2,
    content = $3,
    updated_at = NOW()
where id = $1 and owner_id = $4
returning *;

-- name: CheckIfUserOwned :one

Select EXISTS(
    select 1
    from notes
    where id = $1 and owner_id = $2
);