package main

import (
	"database/sql"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func MigrateHandlerAction(ctx *cli.Context) error {
	config, err := GetConfig(ctx.String("config-file-path"))
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		Debug:            config.Environment != "production",
		AttachStacktrace: true,
		SampleRate:       1.0,
		Release:          version,
		Environment:      config.Environment,
		DebugWriter:      log.Logger,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if config.Environment != "production" {
				log.Debug().Interface("exceptions", event.Exception).Msg(event.Message)
			}

			return event
		},
	})
	if err != nil {
		return fmt.Errorf("initializing Sentry: %w", err)
	}

	conn, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Closing database")
		}
	}()

	migration, err := NewMigration(conn)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	switch ctx.Args().First() {
	case "down":
		err := migration.Down(ctx.Context)
		if err != nil {
			return fmt.Errorf("executing down migration: %w", err)
		}
	case "up":
		fallthrough
	default:
		err := migration.Up(ctx.Context)
		if err != nil {
			return fmt.Errorf("executing up migration: %w", err)
		}
	}

	log.Info().Msg("Migration succeed")

	return nil
}
