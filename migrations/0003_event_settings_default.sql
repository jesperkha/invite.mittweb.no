-- +goose Up
INSERT INTO event_settings (id, title)
VALUES (1, '')
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM event_settings
WHERE id = 1;
