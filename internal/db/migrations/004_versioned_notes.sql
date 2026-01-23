-- +goose Up

CREATE TABLE versioned_notes (
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL,
    note_id UUID references notes(id) ON DELETE CASCADE NOT NULL,
    version_number INT NOT NULL,

    PRIMARY KEY (note_id, version_number)
);

-- +goose Down

drop table versioned_notes;