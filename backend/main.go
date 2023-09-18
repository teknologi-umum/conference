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

	"github.com/flowchartsman/handlebars/v3"
	"github.com/urfave/cli/v2"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

var version string

func App() *cli.App {
	config, err := GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get config")
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		Debug:            config.Environment != "production",
		AttachStacktrace: true,
		SampleRate:       1.0,
		Release:          version,
		Environment:      config.Environment,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Initiating sentry")
	}

	app := &cli.App{
		Name:           "teknum-conf",
		Version:        version,
		Description:    "CLI for working with Teknologi Umum Conference backend",
		DefaultCommand: "server",
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
						UserDomain: userDomain,
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
				Name: "blast-email",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "smtp.hostname",
						Value: "",
						Usage: "SMTP hostname",
					},
					&cli.StringFlag{
						Name:  "smtp.port",
						Value: "",
						Usage: "SMTP port",
					},
					&cli.StringFlag{
						Name:  "smtp.from",
						Value: "admin@localhost",
						Usage: "SMTP sender email",
					},
					&cli.StringFlag{
						Name:  "smtp.password",
						Value: "",
						Usage: "SMTP password",
					},
					&cli.StringFlag{
						Name:     "subject",
						Value:    "",
						Usage:    "Email subject",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "plaintext-body",
						Value:    "",
						Usage:    "Path to plaintext body file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "html-body",
						Value:    "",
						Usage:    "Path to HTML body file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "recipients",
						Value:    "",
						Usage:    "Path to CSV file containing list of emails",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "single-recipient",
						Value:    "",
						Required: false,
					},
				},
				Usage:     "blast-email [subject] [template-plaintext] [template-html-body] [csv-file list destination of emails]",
				ArgsUsage: "[subject] [template-plaintext] [template-html-body] [path-csv-file]",
				Action: func(cCtx *cli.Context) error {
					smtpHostname := cCtx.String("smtp.hostname")
					smtpPort := cCtx.String("smtp.port")
					smtpFrom := cCtx.String("smtp.from")
					smtpPassword := cCtx.String("smtp.password")
					subject := cCtx.String("subject")
					plaintext := cCtx.String("plaintext-body")
					htmlBody := cCtx.String("html-body")
					mailCsv := cCtx.String("recipients")
					singleRecipient := cCtx.String("single-recipient")

					if subject == "" {
						log.Fatal().Msg("Subject is required")
					}
					if plaintext == "" {
						log.Fatal().Msg("Plaintext template is required")
					}
					if htmlBody == "" {
						log.Fatal().Msg("Html template is required")
					}
					if mailCsv == "" && singleRecipient == "" {
						log.Fatal().Msg("Recipient is required")
					}

					plaintextContent, err := os.ReadFile(plaintext)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to read plaintext template")
					}

					plaintextTemplate, err := handlebars.Parse(string(plaintextContent))
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse plaintext template")
					}

					htmlContent, err := os.ReadFile(htmlBody)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to read html template")
					}

					htmlTemplate, err := handlebars.Parse(string(htmlContent))
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse html template")
					}

					var userList []User

					if mailCsv != "" {
						emailList, err := os.ReadFile(mailCsv)
						if err != nil {
							log.Fatal().Err(err).Msg("Failed to read email list")
						}

						userList, err = csvReader(string(emailList))
						if err != nil {
							log.Fatal().Err(err).Msg("Failed to parse email list")
						}
					} else {
						userList = append(userList, User{
							Email: singleRecipient,
						})
					}

					mailSender := NewMailSender(&MailConfiguration{
						SmtpHostname: smtpHostname,
						SmtpPort:     smtpPort,
						SmtpFrom:     smtpFrom,
						SmtpPassword: smtpPassword,
					})

					for _, user := range userList {
						mail := &Mail{
							RecipientName:  user.Name,
							RecipientEmail: user.Email,
							Subject:        subject,
							PlainTextBody:  string(plaintextContent),
							HtmlBody:       string(htmlContent),
						}

						// Execute handlebars template only if user.Name is not empty
						if user.Name != "" {
							mail.PlainTextBody = plaintextTemplate.MustExec(map[string]any{"name": user.Name})
							mail.HtmlBody = htmlTemplate.MustExec(map[string]any{"name": user.Name})
						}

						err := mailSender.Send(mail)
						if err != nil {
							log.Error().Err(err).Msgf("Failed to send email to %s", user.Email)
							continue
						}

						log.Info().Msgf("Sending email to %s", user.Email)
					}
					log.Info().Msg("Blasting email done")
					return nil
				},
			},
		},
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
		},
		Copyright: `   Copyright 2023 Teknologi Umum

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.`,
	}

	return app
}

func main() {
	if err := App().Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}
