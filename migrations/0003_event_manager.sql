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

INSERT INTO event_settings (id, title, description, event_date, event_time, location, contact)
VALUES (1, 'Friend of a Friend', 'Vi inviterer til første Friend of a Friend fest! Alle deltakere skal ta med en +1 som færrest mulig kjenner (helst ingen). Om ingen er tilgjengelig kan du ta med et kjent vesen også. Det blir ulike aktiviteter, mingling, fest, og moro.', 'Torsdag 15. oktober', '18:00', 'Nygårdsgaten 8, 3etg', 'hammerjesp@gmail.com')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS invite_questions (
    id SERIAL PRIMARY KEY,
    label TEXT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0
);

INSERT INTO invite_questions (label, required, sort_order)
SELECT 'Navn på vennen du er +1 til (hvis du er +1)', false, 0
WHERE NOT EXISTS (SELECT 1 FROM invite_questions);

ALTER TABLE registrations ADD COLUMN IF NOT EXISTS answers JSONB NOT NULL DEFAULT '{}'::jsonb;

UPDATE registrations
SET answers = jsonb_build_object('Navn på vennen du er +1 til (hvis du er +1)', friend_name)
WHERE friend_name IS NOT NULL AND friend_name <> '' AND answers = '{}'::jsonb;

ALTER TABLE registrations DROP COLUMN IF EXISTS friend_name;

-- +goose Down
DROP TABLE IF EXISTS invite_questions;
DROP TABLE IF EXISTS event_settings;
ALTER TABLE registrations DROP COLUMN IF EXISTS answers;
ALTER TABLE registrations ADD COLUMN IF NOT EXISTS friend_name TEXT NOT NULL DEFAULT '';
