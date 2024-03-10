package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"conf/mailer"
	"conf/server"
	"conf/ticketing"
	"conf/user"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gocloud.dev/blob"
)

func ServerHandlerAction(ctx *cli.Context) error {
	config, err := GetConfig(ctx.String("config-file-path"))
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		Debug:            config.Environment != "production",
		AttachStacktrace: true,
		SampleRate:       1.0,
		EnableTracing:    true,
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			if ctx.Span.Name == "GET /internal/ping" {
				return 0
			}

			return 0.2
		},
		Release:     version,
		Environment: config.Environment,
		DebugWriter: log.Logger,
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
	defer sentry.Flush(time.Minute)

	pgxRawConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("Parsing connection string configuration")
	}

	pgxConfig := pgxRawConfig.Copy()

	pgxConfig.ConnConfig.Tracer = &PGXTracer{}

	conn, err := pgxpool.NewWithConfig(ctx.Context, pgxConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
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

	mailSender := mailer.NewMailSender(&mailer.MailConfiguration{
		SmtpHostname: config.Mailer.Hostname,
		SmtpPort:     config.Mailer.Port,
		SmtpFrom:     config.Mailer.From,
		SmtpPassword: config.Mailer.Password,
	})

	ticketDomain, err := ticketing.NewTicketDomain(conn, bucket, signaturePrivateKey, signaturePublicKey, mailSender)
	if err != nil {
		return fmt.Errorf("creating ticket domain: %w", err)
	}

	httpServer := server.NewServer(&server.ServerConfig{
		UserDomain:                user.NewUserDomain(conn),
		TicketDomain:              ticketDomain,
		Environment:               config.Environment,
		FeatureRegistrationClosed: config.FeatureFlags.RegistrationClosed,
		ValidateTicketKey:         config.ValidateTicketKey,
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

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("serving http server")
	}

	return nil
}
