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

	"conf/administrator"
	"conf/mailer"
	"conf/nocodb"
	"conf/server"
	"conf/ticketing"
	"conf/user"
	"github.com/getsentry/sentry-go"
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
		Dsn:           "",
		Debug:         config.Environment != "production",
		SampleRate:    1.0,
		EnableTracing: true,
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			if ctx.Span.Name == "GET /api/public/ping" {
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

	database, err := nocodb.NewClient(nocodb.ClientOptions{
		ApiToken:   config.Database.NocoDbApiKey,
		BaseUrl:    config.Database.NocoDbBaseUrl,
		HttpClient: &http.Client{Transport: NewSentryRoundTripper(http.DefaultTransport, nil)},
		Logger:     log.Logger,
	})
	if err != nil {
		return fmt.Errorf("creating database client instance: %w", err)
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

	mailSender := mailer.NewMailSender(&mailer.MailConfiguration{
		SmtpHostname: config.Mailer.Hostname,
		SmtpPort:     config.Mailer.Port,
		SmtpFrom:     config.Mailer.From,
		SmtpPassword: config.Mailer.Password,
	})

	ticketDomain, err := ticketing.NewTicketDomain(database, bucket, signaturePrivateKey, signaturePublicKey, mailSender, config.Database.TicketingTableId)
	if err != nil {
		return fmt.Errorf("creating ticket domain: %w", err)
	}

	userDomain, err := user.NewUserDomain(database, config.Database.UserTableId)
	if err != nil {
		return fmt.Errorf("creating user domain: %w", err)
	}

	administratorDomain, err := administrator.NewAdministratorDomain(config.AdministratorUserMapping)
	if err != nil {
		return fmt.Errorf("creating administrator domain: %w", err)
	}

	httpServer, err := server.NewServer(&server.ServerConfig{
		UserDomain:          userDomain,
		TicketDomain:        ticketDomain,
		AdministratorDomain: administratorDomain,
		FeatureFlag:         &config.FeatureFlags,
		MailSender:          mailSender,
		Environment:         config.Environment,
		ValidateTicketKey:   config.ValidateTicketKey,
		Hostname:            "",
		Port:                config.Port,
	})
	if err != nil {
		return fmt.Errorf("creating http server: %w", err)
	}

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
