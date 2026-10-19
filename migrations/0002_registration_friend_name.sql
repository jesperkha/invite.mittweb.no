-- +goose Up
ALTER TABLE registrations
ADD COLUMN friend_name TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE registrations
DROP COLUMN friend_name;
