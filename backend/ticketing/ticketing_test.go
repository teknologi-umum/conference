package ticketing_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"conf/mailer"
	"conf/nocodb"
	"conf/nocodb/nocodbmock"
	"conf/ticketing"
	"conf/user"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

var database *nocodb.Client
var bucket *blob.Bucket
var mailSender *mailer.Mailer

func TestMain(m *testing.M) {
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

	nocodbMockServer, err := nocodbmock.NewNocoDBMockServer()
	if err != nil {
		log.Fatal().Err(err).Msg("creating nocodb mock server")
		return
	}

	database, err = nocodb.NewClient(nocodb.ClientOptions{
		ApiToken:   "testing",
		BaseUrl:    nocodbMockServer.URL,
		HttpClient: nocodbMockServer.Client(),
		Logger:     log.Logger,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("creating nocodb client")
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
	nocodbMockServer.Close()

	os.Exit(exitCode)
}

func TestNewTicketDomain(t *testing.T) {
	// Create mock dependencies.
	db := &nocodb.Client{}
	bucket := &blob.Bucket{}
	privateKey := ed25519.PrivateKey{}
	publicKey := ed25519.PublicKey{}
	mailSender := &mailer.Mailer{}

	// Group the tests with t.Run().
	t.Run("all dependencies set", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, privateKey, publicKey, mailSender)
		if err != nil {
			t.Errorf("NewTicketDomain failed: %v", err)
		}
		if ticketDomain == nil {
			t.Error("NewTicketDomain returned nil ticketDomain")
		}
	})

	t.Run("nil database", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(nil, bucket, privateKey, publicKey, mailSender)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil database")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil database")
		}
	})

	t.Run("nil bucket", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, nil, privateKey, publicKey, mailSender)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil bucket")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil bucket")
		}
	})

	t.Run("nil private key", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, nil, publicKey, mailSender)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil private key")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil private key")
		}
	})

	t.Run("nil public key", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, privateKey, nil, mailSender)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil public key")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil public key")
		}
	})

	t.Run("nil mailSender", func(t *testing.T) {
		ticketDomain, err := ticketing.NewTicketDomain(db, bucket, privateKey, publicKey, nil)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil mailSender")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil mailSender")
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

	userDomain, err := user.NewUserDomain(database, "testing")
	if err != nil {
		t.Fatalf("creating user domain instance: %s", err.Error())
	}

	t.Run("Invalid photo", func(t *testing.T) {
		err := ticketDomain.StorePaymentReceipt(context.Background(), user.User{}, nil, "")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		var validationError *ticketing.ValidationError
		if errors.As(err, &validationError) {
			if len(validationError.Errors) != 2 {
				t.Errorf("expecting two errors, got %d", len(validationError.Errors))
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

		err = ticketDomain.StorePaymentReceipt(ctx, user.User{Email: email}, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
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

		user := user.User{
			Email: email,
		}

		// First attempt
		err = ticketDomain.StorePaymentReceipt(ctx, user, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		// Second attempt, should not return error
		err = ticketDomain.StorePaymentReceipt(ctx, user, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err != nil {
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

	userDomain, err := user.NewUserDomain(database, "testing")
	if err != nil {
		t.Fatalf("creating user domain instance: %s", err.Error())
	}

	t.Run("Email not found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		_, err := ticketDomain.ValidatePaymentReceipt(ctx, user.User{Email: "not-found@example.com"})
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

		user := user.User{Email: email}

		err = ticketDomain.StorePaymentReceipt(ctx, user, strings.NewReader("Hello world! This is not a photo. Yet this will be a text file."), "text/plain")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		sum, err := ticketDomain.ValidatePaymentReceipt(ctx, user)
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

	t.Run("Happy scenario", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := ticketDomain.VerifyIsStudent(ctx, user.User{Email: "aji@test.com"})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}
func TestNullTicketing_MarshalJSON(t *testing.T) {
	n := ticketing.NullTicketing{
		Id:               sql.NullInt64{Valid: true, Int64: 123},
		Email:            sql.NullString{},
		ReceiptPhotoPath: sql.NullString{},
		Paid:             sql.NullBool{},
		SHA256Sum:        sql.NullString{},
		Used:             sql.NullBool{},
		CreatedAt:        sql.NullTime{Valid: false},
		UpdatedAt:        sql.NullTime{Valid: true, Time: time.Time{}},
	}

	out, err := json.Marshal(n)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	expect := `{"Id":123,"UpdatedAt":"0001-01-01T00:00:00Z"}`
	if string(out) != expect {
		t.Errorf("expecting %s, got %s", expect, string(out))
	}
}
