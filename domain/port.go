package domain

import (
	"context"
	"encoding/json"
	"time"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type RegisterRepository interface {
	Register(ctx context.Context, name, email string, answers map[string]string) error
	GetRegistrations(ctx context.Context) ([]Registration, error)
	DeleteRegistration(ctx context.Context, id int) error
	DeleteRegistrations(ctx context.Context) error
	GetEvent(ctx context.Context) (Event, error)
	UpdateEvent(ctx context.Context, event Event) error
}

type Registration struct {
	ID        int             `db:"id" json:"id"`
	Name      string          `db:"name" json:"name"`
	Email     string          `db:"email" json:"email"`
	Answers   json.RawMessage `db:"answers" json:"answers"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}

type Question struct {
	ID       int    `json:"id"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
}

type Event struct {
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Date        string     `db:"event_date" json:"date"`
	Time        string     `db:"event_time" json:"time"`
	Location    string     `db:"location" json:"location"`
	Contact     string     `db:"contact" json:"contact"`
	ButtonColor string     `db:"button_color" json:"button_color"`
	Background  string     `db:"background" json:"background"`
	Questions   []Question `json:"questions"`
}
