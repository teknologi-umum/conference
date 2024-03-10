package ticketing_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"conf/mailer"
	"conf/ticketing"
	"conf/user"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

var database *pgxpool.Pool
var bucket *blob.Bucket
var mailSender *mailer.Mailer

func TestMain(m *testing.M) {
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		databaseUrl = "postgres://postgres:password@localhost:5432/conf?sslmode=disable"
	}

	tempDir, err := os.MkdirTemp(os.TempDir(), "teknologi-umum-conference")
	if err != nil {
		log.Fatal().Err(err).Msg("creating temporary directory")
		return
	}

	blobUrl, ok := os.LookupEnv("BLOB_URL")
	if !ok {
		blobUrl = "file://" + tempDir
	}

	smtpHostname, ok := os.LookupEnv("SMTP_HOSTNAME")
	if !ok {
		smtpHostname = "localhost"
	}
	smtpPort, ok := os.LookupEnv("SMTP_PORT")
	if !ok {
		smtpPort = "1025"
	}
	smtpFrom, ok := os.LookupEnv("SMTP_FROM")
	if !ok {
		smtpFrom = ""
	}
	smtpPassword, ok := os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		smtpPassword = ""
	}

	_ = sentry.Init(sentry.ClientOptions{})

	database, err = pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("creating pgx pool instance")
		return
	}

	bucket, err = blob.OpenBucket(context.Background(), blobUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("creating bucket instance")
		return
	}

	mailSender = mailer.NewMailSender(&mailer.MailConfiguration{
		SmtpHostname: smtpHostname,
		SmtpPort:     smtpPort,
		SmtpFrom:     smtpFrom,
		SmtpPassword: smtpPassword,
	})

	exitCode := m.Run()

	_ = os.RemoveAll(tempDir)
	_ = bucket.Close()
	database.Close()

	os.Exit(exitCode)
}

func TestNewTicketDomain(t *testing.T) {
	// Create mock dependencies.
	db := &pgxpool.Pool{}
	bucket := &blob.Bucket{}
	privateKey := ed25519.PrivateKey{}
	publicKey := ed25519.PublicKey{}
	mailer := &mailer.Mailer{}

	// Group the tests with t.Run().
	t.Run("all dependencies set", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, privateKey, publicKey, mailer)
		if err != nil {
			t.Errorf("NewTicketDomain failed: %v", err)
		}
		if ticketDomain == nil {
			t.Error("NewTicketDomain returned nil ticketDomain")
		}
	})

	t.Run("nil database", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(nil, bucket, privateKey, publicKey, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil database")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil database")
		}
	})

	t.Run("nil bucket", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, nil, privateKey, publicKey, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil bucket")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil bucket")
		}
	})

	t.Run("nil private key", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, nil, publicKey, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil private key")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil private key")
		}
	})

	t.Run("nil public key", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, privateKey, nil, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil public key")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil public key")
		}
	})

	t.Run("nil mailer", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, privateKey, publicKey, nil)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil mailer")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil mailer")
		}
	})
}

func TestTicketDomain_StorePaymentReceipt(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating new ed25519 key: %s", err.Error())
		return
	}

	ticketDomain, err := ticketing.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender)
	if err != nil {
		t.Fatalf("creating a ticket domain instance: %s", err.Error())
	}

	userDomain := user.NewUserDomain(database)

	t.Run("Invalid Email and photo", func(t *testing.T) {
		err := ticketDomain.StorePaymentReceipt(context.Background(), "", nil, "")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		var validationError *ticketing.ValidationError
		if errors.As(err, &validationError) {
			if len(validationError.Errors) != 3 {
				t.Errorf("expecting three errors, got %d", len(validationError.Errors))
			}
		}
	})

	t.Run("Happy scenario", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		email := "johndoe+happy@example.com"
		err := userDomain.CreateParticipant(ctx, user.CreateParticipantRequest{
			Name:  "John Doe",
			Email: email,
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		err = ticketDomain.StorePaymentReceipt(ctx, email, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

	t.Run("Update data if email already exists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		email := "johndoe+happy@example.com"
		err := userDomain.CreateParticipant(ctx, user.CreateParticipantRequest{
			Name:  "John Doe",
			Email: email,
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// First attempt
		err = ticketDomain.StorePaymentReceipt(ctx, email, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// Second attempt, should not return error
		err = ticketDomain.StorePaymentReceipt(ctx, email, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

	t.Run("User email not found, should return ErrUserEmailNotFound", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := ticketDomain.StorePaymentReceipt(ctx, "johndoe+not+found@example.com", strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		if err != nil && !errors.Is(err, ticketing.ErrUserEmailNotFound) {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}

func TestTicketDomain_ValidatePaymentReceipt(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating new ed25519 key: %s", err.Error())
		return
	}

	ticketDomain, err := ticketing.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender)
	if err != nil {
		t.Fatalf("creating a ticket domain instance: %s", err.Error())
	}

	userDomain := user.NewUserDomain(database)

	t.Run("Invalid email", func(t *testing.T) {
		_, err := ticketDomain.ValidatePaymentReceipt(context.Background(), "")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		var validationError *ticketing.ValidationError
		if errors.As(err, &validationError) {
			if len(validationError.Errors) != 1 {
				t.Errorf("expecting one error, got %d", len(validationError.Errors))
			}
		}
	})

	t.Run("Email not found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		_, err := ticketDomain.ValidatePaymentReceipt(ctx, "not-found@example.com")
		if err == nil {
			t.Error("expecting an error, got nil")
		}

		if !errors.Is(err, ticketing.ErrInvalidTicket) {
			t.Errorf("expecting an error of ErrInvalidTicket, instead got %s", err.Error())
		}
	})

	t.Run("Happy scenario", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		email := "johndoe+happy@example.com"
		err := userDomain.CreateParticipant(ctx, user.CreateParticipantRequest{
			Name:  "John Doe",
			Email: email,
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		err = ticketDomain.StorePaymentReceipt(ctx, email, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		sum, err := ticketDomain.ValidatePaymentReceipt(ctx, email)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if sum == "" {
			t.Error("expecting sum to have value, got empty string")
		}
	})
}

func TestTicketDomain_VerifyIsStudent(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating new ed25519 key: %s", err.Error())
		return
	}
	ticketDomain, err := ticketing.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender)
	if err != nil {
		t.Fatalf("creating a ticket domain instance: %s", err.Error())
	}
	t.Run("Invalid email", func(t *testing.T) {
		err := ticketDomain.VerifyIsStudent(context.Background(), "")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		var validationError *ticketing.ValidationError
		if errors.As(err, &validationError) {
			if len(validationError.Errors) != 1 {
				t.Errorf("expecting one error, got %d", len(validationError.Errors))
			}
		}
	})

	t.Run("Happy scenario", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := ticketDomain.VerifyIsStudent(ctx, "aji@test.com")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}
