package repo

import (
	"context"
	"singel-sondag/domain"
)

func (db *Database) GetEvent(ctx context.Context) (domain.Event, error) {
	var event domain.Event
	err := db.GetContext(ctx, &event, `SELECT title, description, event_date, event_time, location, contact, button_color, background FROM event_settings WHERE id = 1`)
	if err != nil {
		return event, err
	}
	event.Questions = []domain.Question{}
	err = db.SelectContext(ctx, &event.Questions, `SELECT id, label, required FROM invite_questions ORDER BY sort_order, id`)
	return event, err
}

func (db *Database) UpdateEvent(ctx context.Context, event domain.Event) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `UPDATE event_settings SET title=$1, description=$2, event_date=$3, event_time=$4, location=$5, contact=$6, button_color=$7, background=$8 WHERE id=1`,
		event.Title, event.Description, event.Date, event.Time, event.Location, event.Contact, event.ButtonColor, event.Background)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM invite_questions`)
	if err != nil {
		return err
	}
	for i, question := range event.Questions {
		if question.Label == "" {
			continue
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO invite_questions (label, required, sort_order) VALUES ($1, $2, $3)`, question.Label, question.Required, i); err != nil {
			return err
		}
	}
	return tx.Commit()
}
