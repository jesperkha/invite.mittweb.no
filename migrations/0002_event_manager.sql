-- +goose Up
CREATE TABLE IF NOT EXISTS event_settings (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    event_date TEXT NOT NULL DEFAULT '',
    event_time TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    contact TEXT NOT NULL DEFAULT '',
    button_color TEXT NOT NULL DEFAULT '#F28482',
    background TEXT NOT NULL DEFAULT '#F7EDE2'
);

CREATE TABLE IF NOT EXISTS invite_questions (
    id SERIAL PRIMARY KEY,
    label TEXT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0
);

ALTER TABLE registrations ADD COLUMN IF NOT EXISTS answers JSONB NOT NULL DEFAULT '{}'::jsonb;

-- +goose Down
DROP TABLE IF EXISTS invite_questions;
DROP TABLE IF EXISTS event_settings;
ALTER TABLE registrations DROP COLUMN IF EXISTS answers;
