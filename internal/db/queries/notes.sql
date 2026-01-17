-- name: GetAllNotes :many

Select 
     id ,
    title,
    content,
    updated_at
From notes;

-- name: CreateNewNote :one

Insert into notes (title, content , updated_at)
Values ($1,$2 , NOW())
Returning *;

-- name: GetANote :one

Select 
    id,
    title,
    content,
    updated_at
  from notes
where id = $1;

-- name: DeleteNote :one
Delete from notes
where id = $1
returning *;

-- name: UpdateNote :one

update notes
Set 
    title = $2,
    content = $3,
    updated_at = NOW()
where id = $1
returning *;