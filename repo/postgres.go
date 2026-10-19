package repo

import (
	"context"
	"fmt"
	"singel-sondag/domain"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
	logger domain.Logger
}

func New(logger domain.Logger, databaseUrl string) (*Database, error) {
	db, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	logger.Info(context.Background(), "connected to database")

	return &Database{
		db,
		logger,
	}, nil
}

func (db *Database) gooseUp(path string, tableName string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	goose.SetTableName(tableName)

	if err := goose.Up(db.DB.DB, path); err != nil {
		return fmt.Errorf("failed to run goose up: %w", err)
	}

	return nil
}

func (db *Database) Migrate(path string) error {
	err := db.gooseUp(path, "goose_db_version")
	db.logger.Info(context.Background(), "finished migrations")
	return err
}

func (db *Database) Seed(path string) error {
	err := db.gooseUp(path, "goose_seed_version")
	db.logger.Info(context.Background(), "finished seeding")
	return err
}

func (db *Database) Close() error {
	return db.DB.Close()
}

// Implement Connection
func (db *Database) Rollback() error {
	return nil
}
