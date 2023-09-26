package main

import (
	"context"
	"database/sql"
	"encoding/hex"
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
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

var version string

func App() *cli.App {
	app := &cli.App{
		Name:           "teknum-conf",
		Version:        version,
		Description:    "CLI for working with Teknologi Umum Conference backend",
		DefaultCommand: "server",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-file-path",
				EnvVars: []string{"CONFIGURATION_FILE"},
			},
		},
		Commands: []*cli.Command{
			{
				Name: "server",
				Action: func(cCtx *cli.Context) error {
					config, err := GetConfig(cCtx.String("config-file-path"))
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

					conn, err := pgxpool.New(
						context.Background(),
						fmt.Sprintf(
							"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
							config.Database.User,
							config.Database.Password,
							config.Database.Host,
							config.Database.Port,
							config.Database.Name,
						),
					)
					if err != nil {
						return fmt.Errorf("failed connect to database: %w", err)
					}
					defer conn.Close()

					bucket, err := blob.OpenBucket(context.Background(), config.BlobUrl)
					if err != nil {
						return fmt.Errorf("opening bucket: %w", err)
					}
					defer func() {
						err := bucket.Close()
						if err != nil {
							log.Warn().Err(err).Msg("Closing bucket")
						}
					}()

					signaturePrivateKey, err := hex.DecodeString(config.Signature.PrivateKey)
					if err != nil {
						return fmt.Errorf("invalid signature private key: %w", err)
					}

					signaturePublicKey, err := hex.DecodeString(config.Signature.PublicKey)
					if err != nil {
						return fmt.Errorf("invalid signature public key: %w", err)
					}

					mailer := NewMailSender(&MailConfiguration{
						SmtpHostname: config.Mailer.Hostname,
						SmtpPort:     config.Mailer.Port,
						SmtpFrom:     config.Mailer.From,
						SmtpPassword: config.Mailer.Password,
					})

					ticketDomain, err := NewTicketDomain(conn, bucket, signaturePrivateKey, signaturePublicKey, mailer)
					if err != nil {
						return fmt.Errorf("creating ticket domain: %w", err)
					}

					httpServer := NewServer(&ServerConfig{
						UserDomain:   NewUserDomain(conn),
						TicketDomain: ticketDomain,
						Environment:  config.Environment,
					})

					exitSig := make(chan os.Signal, 1)
					signal.Notify(exitSig, os.Interrupt, os.Kill)

					go func() {
						<-exitSig
						ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
						defer cancel()

						if err := httpServer.Shutdown(ctx); err != nil {
							log.Error().Err(err).Msg("failed to shutdown server")
						}
					}()

					if err := httpServer.Start(net.JoinHostPort("", config.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
						log.Fatal().Err(err).Msg("failed to start server")
					}

					return nil
				},
			},
			{
				Name: "migrate",
				Action: func(cCtx *cli.Context) error {
					config, err := GetConfig(cCtx.String("config-file-path"))
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

					switch cCtx.Args().First() {
					case "down":
						err := migration.Down(cCtx.Context)
						if err != nil {
							return fmt.Errorf("executing down migration: %w", err)
						}
					case "up":
						fallthrough
					default:
						err := migration.Up(cCtx.Context)
						if err != nil {
							return fmt.Errorf("executing up migration: %w", err)
						}
					}

					log.Info().Msg("Migration succeed")

					return nil
				},
			},
			{
				Name: "dump-attendees",
				Action: func(cCtx *cli.Context) error {
					config, err := GetConfig(cCtx.String("config-file-path"))
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

					conn, err := pgxpool.New(
						cCtx.Context,
						fmt.Sprintf(
							"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
							config.Database.User,
							config.Database.Password,
							config.Database.Host,
							config.Database.Port,
							config.Database.Name,
						),
					)
					if err != nil {
						log.Fatal().Err(err).Msg("failed to connect to database")
					}
					defer conn.Close()

					userDomain := NewUserDomain(conn)

					return userDomain.ExportUnprocessedUsersAsCSV(cCtx.Context)
				},
			},
			{
				Name: "blast-email",
				Flags: []cli.Flag{
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
					config, err := GetConfig(cCtx.String("config-file-path"))
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
						log.Fatal().Err(err).Msg("failed to read plaintext template")
					}

					plaintextTemplate, err := handlebars.Parse(string(plaintextContent))
					if err != nil {
						log.Fatal().Err(err).Msg("failed to parse plaintext template")
					}

					htmlContent, err := os.ReadFile(htmlBody)
					if err != nil {
						log.Fatal().Err(err).Msg("failed to read html template")
					}

					htmlTemplate, err := handlebars.Parse(string(htmlContent))
					if err != nil {
						log.Fatal().Err(err).Msg("failed to parse html template")
					}

					var userList []User

					if mailCsv != "" {
						emailList, err := os.ReadFile(mailCsv)
						if err != nil {
							log.Fatal().Err(err).Msg("failed to read email list")
						}

						userList, err = csvReader(string(emailList), true)
						if err != nil {
							log.Fatal().Err(err).Msg("failed to parse email list")
						}
					} else {
						userList = append(userList, User{
							Email: singleRecipient,
						})
					}

					mailSender := NewMailSender(&MailConfiguration{
						SmtpHostname: config.Mailer.Hostname,
						SmtpPort:     config.Mailer.Port,
						SmtpFrom:     config.Mailer.From,
						SmtpPassword: config.Mailer.Password,
					})

					for _, user := range userList {
						mail := &Mail{
							RecipientName:  user.Name,
							RecipientEmail: user.Email,
							Subject:        subject,
							PlainTextBody:  string(plaintextContent),
							HtmlBody:       string(htmlContent),
						}


						// Parse email template information
						emailTemplate := map[string]any{
							"ticketPrice":                         config.EmailTemplate.TicketPrice,
							"ticketStudentCollegePrice":           config.EmailTemplate.TicketStudentCollegePrice,
							"ticketStudentHighSchoolPrice":        config.EmailTemplate.TicketStudentHighSchoolPrice,
							"ticketStudentCollegeDiscount":        config.EmailTemplate.TicketStudentCollegeDiscount,
							"ticketStudentHighSchoolDiscount":     config.EmailTemplate.TicketStudentHighSchoolDiscount,
							"percentageStudentCollegeDiscount":    config.EmailTemplate.PercentageStudentCollegeDiscount,
							"percentageStudentHighSchoolDiscount": config.EmailTemplate.PercentageStudentHighSchoolDiscount,
							"conferenceEmail":                     config.EmailTemplate.ConferenceEmail,
							"bankAccounts":                        config.EmailTemplate.BankAccounts,
						}
						// Execute handlebars template only if user.Name is not empty
						if user.Name != "" {
							emailTemplate["name"] = user.Name
						}


						mail.PlainTextBody = plaintextTemplate.MustExec(emailTemplate)
						mail.HtmlBody = htmlTemplate.MustExec(emailTemplate)

						err := mailSender.Send(mail)
						if err != nil {
							log.Error().Err(err).Msgf("failed to send email to %s", user.Email)
							continue
						}

						log.Info().Msgf("Sending email to %s", user.Email)
					}
					log.Info().Msg("Blasting email done")
					return nil
				},
			},
			{
				Name:      "participants",
				Usage:     "participants [is_processed]",
				ArgsUsage: "[is_processed]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "is_processed",
						Value: false,
						Usage: "Is processed",
					},
				},
				Action: func(cCtx *cli.Context) error {
					config, err := GetConfig(cCtx.String("config-file-path"))
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

					isProcessedStr := cCtx.Bool("is_processed")

					conn, err := pgxpool.New(
						cCtx.Context,
						fmt.Sprintf(
							"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
							config.Database.User,
							config.Database.Password,
							config.Database.Host,
							config.Database.Port,
							config.Database.Name,
						),
					)
					if err != nil {
						return err
					}
					defer conn.Close()

					userDomain := NewUserDomain(conn)
					users, err := userDomain.GetUsers(cCtx.Context, UserFilterRequest{Type: TypeParticipant, IsProcessed: isProcessedStr})
					if err != nil {
						return err
					}

					log.Info().Msg("List of participants")
					for _, user := range users {
						log.Info().Msgf("%s - %s", user.Name, user.Email)
					}

					return nil
				},
			},
			{
				Name:      "student-verification",
				Usage:     "student-verification [path-csv-file]",
				ArgsUsage: "[path-csv-file]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "bulk-verification",
						Value:    "",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "single-verification",
						Value:    "",
						Required: false,
					},
				},
				Action: func(cCtx *cli.Context) error {
					config, err := GetConfig(cCtx.String("config-file-path"))
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

					bulkVerification := cCtx.String("bulk-verification")
					singleVerification := cCtx.String("single-verification")

					if bulkVerification == "" && singleVerification == "" {
						return fmt.Errorf("requires `--bulk-verification` or `--single-verification` flag")
					}

					var students []User
					if bulkVerification != "" {
						emailList, err := os.ReadFile(bulkVerification)
						if err != nil {
							log.Fatal().Err(err).Msg("failed to read email list")
						}

						students, err = csvReader(string(emailList), false)
						if err != nil {
							log.Fatal().Err(err).Msg("failed to parse email list")
						}
					} else {
						students = append(students, User{
							Email: singleVerification,
						})
					}

					bucket, err := blob.OpenBucket(context.Background(), config.BlobUrl)
					if err != nil {
						return fmt.Errorf("opening bucket: %w", err)
					}
					defer func() {
						err := bucket.Close()
						if err != nil {
							log.Warn().Err(err).Msg("Closing bucket")
						}
					}()

					signaturePrivateKey, err := hex.DecodeString(config.Signature.PrivateKey)
					if err != nil {
						return fmt.Errorf("invalid signature private key: %w", err)
					}

					signaturePublicKey, err := hex.DecodeString(config.Signature.PublicKey)
					if err != nil {
						return fmt.Errorf("invalid signature public key: %w", err)
					}

					mailer := NewMailSender(&MailConfiguration{
						SmtpHostname: config.Mailer.Hostname,
						SmtpPort:     config.Mailer.Port,
						SmtpFrom:     config.Mailer.From,
						SmtpPassword: config.Mailer.Password,
					})

					conn, err := pgxpool.New(
						cCtx.Context,
						fmt.Sprintf(
							"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
							config.Database.User,
							config.Database.Password,
							config.Database.Host,
							config.Database.Port,
							config.Database.Name,
						),
					)
					if err != nil {
						return fmt.Errorf("failed to connect to database: %w", err)
					}

					ticketDomain, err := NewTicketDomain(conn, bucket, signaturePrivateKey, signaturePublicKey, mailer)
					if err != nil {
						return fmt.Errorf("creating a ticket domain instance: %s", err.Error())
					}

					for _, student := range students {
						err := ticketDomain.VerifyIsStudent(cCtx.Context, student.Email)
						if err != nil {
							log.Error().Err(err).Msgf("failed to verify student %s", student.Email)
							continue
						}

						log.Info().Msgf("Verified student %s", student.Email)
					}

					return nil
				},
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
		log.Fatal().Err(err).Msg("failed to run app")
	}
}
