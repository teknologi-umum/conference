package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/rs/zerolog/log"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config, err := GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get config")
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		Debug:            config.Environment != "production",
		AttachStacktrace: false,
		SampleRate:       1.0,
		Environment:      config.Environment,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Initiating sentry")
		return
	}

	app := &cli.App{
		Name:  "teknum-conf",
		Usage: "say a greeting",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       config.DBName,
				Usage:       "db name",
				Destination: &config.DBName,
			},
			&cli.StringFlag{
				Name:        "db-user",
				Value:       config.DBUser,
				Usage:       "db user",
				Destination: &config.DBUser,
			},
			&cli.StringFlag{
				Name:        "db-password",
				Value:       config.DBPassword,
				Usage:       "db password",
				Destination: &config.DBPassword,
			},
			&cli.StringFlag{
				Name:        "db-host",
				Value:       config.DBHost,
				Usage:       "db host",
				Destination: &config.DBHost,
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       config.Port,
				Usage:       "port",
				Destination: &config.Port,
			},
			&cli.StringFlag{
				Name:        "smtp-hostname",
				Value:       config.SmtpHostname,
				Usage:       "smtp hostname",
				Destination: &config.SmtpHostname,
			},
			&cli.StringFlag{
				Name:        "smtp-port",
				Value:       config.SmtpPort,
				Usage:       "smtp port",
				Destination: &config.SmtpPort,
			},
			&cli.StringFlag{
				Name:        "smtp-from",
				Value:       config.SmtpFrom,
				Usage:       "smtp from",
				Destination: &config.SmtpFrom,
			},
			&cli.StringFlag{
				Name:        "smtp-password",
				Value:       config.SmtpPassword,
				Usage:       "smtp password",
				Destination: &config.SmtpPassword,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "server",
				Action: func(cCtx *cli.Context) error {
					conn, err := pgxpool.New(
						context.Background(),
						fmt.Sprintf(
							"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
							config.DBUser,
							config.DBPassword,
							config.DBHost,
							config.DBPort,
							config.DBName,
						),
					)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to connect to database")
					}
					defer conn.Close()

					userDomain := NewUserDomain(conn)

					e := NewServer(&ServerConfig{
						userDomain: userDomain,
					})

					exitSig := make(chan os.Signal, 1)
					signal.Notify(exitSig, os.Interrupt, os.Kill)

					go func() {
						<-exitSig
						ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
						defer cancel()

						if err := e.Shutdown(ctx); err != nil {
							log.Error().Err(err).Msg("Failed to shutdown server")
						}
					}()

					if err := e.Start(net.JoinHostPort("", config.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
						log.Fatal().Err(err).Msg("Failed to start server")
					}

					return nil
				},
			},
			{
				Name: "migrate",
				Action: func(cCtx *cli.Context) error {
					conn, err := sql.Open(
						"pgx",
						fmt.Sprintf(
							"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
							config.DBUser,
							config.DBPassword,
							config.DBHost,
							config.DBPort,
							config.DBName,
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
					switch cCtx.Args().First() {
					case "down":
						err := migration.Down(cCtx.Context)
						if err != nil {
							return fmt.Errorf("Executing down migration: %w", err)
						}
					case "up":
						fallthrough
					default:
						err := migration.Up(cCtx.Context)
						if err != nil {
							return fmt.Errorf("Executing up migration: %w", err)
						}
					}

					log.Info().Msg("Migration succeed")

					return nil
				},
			},
			{
				Name: "dump-attendees",
				Action: func(cCtx *cli.Context) error {
					conn, err := pgxpool.New(
						context.Background(),
						fmt.Sprintf(
							"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
							config.DBUser,
							config.DBPassword,
							config.DBHost,
							config.DBPort,
							config.DBName,
						),
					)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to connect to database")
					}
					defer conn.Close()

					userDomain := NewUserDomain(conn)

					err = userDomain.ExportUnprocessedUsersAsCSV(cCtx.Context)
					return err
				},
			},
			{
				Name:      "blast-email",
				Usage:     "blast-email [template] [file list destination of emails]",
				ArgsUsage: "[template] [file list destination of emails]",
				Action: func(cCtx *cli.Context) error {
					templateArg := cCtx.Args().Get(0)
					emailList := cCtx.Args().Tail()

					if templateArg == "" {
						log.Fatal().Msg("Template is required")
					}
					if len(emailList) == 0 {
						log.Fatal().Msg("Email list is required for blasting email minimum 1 email")
					}

					htmlContent, err := os.ReadFile(templateArg)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to read template")
					}

					mailSender := NewMailSender(&MailConfiguration{
						SmtpHostname: config.SmtpHostname,
						SmtpPort:     config.SmtpPort,
						SmtpFrom:     config.SmtpFrom,
						SmtpPassword: config.SmtpPassword,
					})
					for _, email := range emailList {
						mail := &Mail{
							RecipientName:  email,
							RecipientEmail: email,
							Subject:        "TeknumConf - Attendee Waitlist",
							PlainTextBody: fmt.Sprintf(`Hello, %s!
					Thank you for participating on TeknumConf, due to the limited seating quota, you are on a waitlist.
					Not to worry, you will receive an email from us regarding the seating in about 7 days.
					Please do contact us if you didn't receive any email by then.`, email),
							HtmlBody: string(htmlContent),
						}
						err := mailSender.Send(mail)
						if err != nil {
							log.Error().Err(err).Msgf("Failed to send email to %s", email)
							continue
						}
						log.Info().Msgf("Sending email to %s", email)
					}

					log.Info().Msg("Blasting email done")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}
