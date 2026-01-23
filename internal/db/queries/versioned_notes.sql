-- name: CreateAVersion :one

Insert into versioned_notes(note_id,version_number,title,content,updated_at)
Values(
  $1,
  $2,
  $3,
  $4,
  NOW()
) 
returning note_id,version_number,updated_at;

-- name: GetCurrentHighestVersion :one

Select COALESCE(MAX(version_number) , 0)::INT
From versioned_notes 
where note_id = $1;

-- name: GetAllVersions :many

Select version_number , created_at
From versioned_notes where note_id = $1;

-- name: GetNoteVersion :one

Select * From versioned_notes
where note_id = $1 and version_number = $2;