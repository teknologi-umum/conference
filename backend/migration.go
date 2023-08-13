package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Migration struct {
	db *sql.DB
}

func NewMigration(db *sql.DB) (*Migration, error) {
	if db == nil {
		return &Migration{}, errors.New("db is nil")
	}

	goose.SetBaseFS(embedMigrations)

	goose.SetLogger(&conformedLogger{})

	if err := goose.SetDialect("postgres"); err != nil {
		return &Migration{}, err
	}

	return &Migration{db: db}, nil
}

func (m *Migration) Up(ctx context.Context) (err error) {
	return goose.UpContext(ctx, m.db, "migrations")
}

func (m *Migration) Down(ctx context.Context) error {
	return goose.DownContext(ctx, m.db, "migrations")
}
