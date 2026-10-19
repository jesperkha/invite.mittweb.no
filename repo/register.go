package repo

import (
	"context"
	"encoding/json"
	"singel-sondag/domain"
)

func (db *Database) Register(ctx context.Context, name string, email string, answers map[string]string) error {
	db.logger.Info(ctx, "registering", "name", name, "email", email)
	answerJSON, err := json.Marshal(answers)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "INSERT INTO registrations (name, email, answers) VALUES ($1, $2, $3)", name, email, answerJSON)
	if err != nil {
		db.logger.Error(ctx, "failed to register", "err", err)
		return err
	}
	return nil
}

func (db *Database) DeleteRegistration(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM registrations WHERE id = $1", id)
	return err
}

func (db *Database) DeleteRegistrations(ctx context.Context) error {
	_, err := db.ExecContext(ctx, "DELETE FROM registrations")
	return err
}

func (db *Database) GetRegistrations(ctx context.Context) ([]domain.Registration, error) {
	db.logger.Info(ctx, "getting registrations")
	var registrations []domain.Registration
	err := db.SelectContext(ctx, &registrations, "SELECT * FROM registrations")
	if err != nil {
		db.logger.Error(ctx, "failed to get registrations", "err", err)
		return nil, err
	}
	return registrations, nil
}
