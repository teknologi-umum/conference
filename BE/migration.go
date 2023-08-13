package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Migration struct {
	DB *pgxpool.Pool
}

func MigrationNew(DB *pgxpool.Pool) *Migration {
	return &Migration{DB: DB}
}

func (m *Migration) Migrate(ctx context.Context) (err error) {
	err = m.DB.Ping(ctx)
	if err != nil {
		err = fmt.Errorf("failed to ping database: %w", err)
		return
	}

	tx, err := m.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %w", err)
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Err(err).Msg("failed to rollback transaction")
			}

		} else {
			log.Info().Msg("migration success")
		}
	}()

	for i, queries := range V1() {
		_, err = tx.Exec(ctx, queries)
		if err != nil {
			err = fmt.Errorf("failed to execute migration at %d: with error %w", i, err)
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return
	}
	return
}

func V1() []string {
	return []string{
		`CREATE TYPE user_type AS ENUM ('participant', 'speaker');`,
		`CREATE TABLE IF NOT EXISTS users (
    			id SERIAL PRIMARY KEY,
    			name VARCHAR(255) NOT NULL,
    			email VARCHAR(255) NOT NULL,
				type user_type,
    			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}
}
