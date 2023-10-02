package main_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gocloud.dev/blob"

	main "conf"
)

func TestNewTicketDomain(t *testing.T) {
	// Create mock dependencies.
	db := &pgxpool.Pool{}
	bucket := &blob.Bucket{}
	privateKey := ed25519.PrivateKey{}
	publicKey := ed25519.PublicKey{}
	mailer := &main.Mailer{}

	// Group the tests with t.Run().
	t.Run("all dependencies set", func(t *testing.T) {
		ticketDomain, err := main.NewTicketDomain(db, bucket, privateKey, publicKey, mailer)
		if err != nil {
			t.Errorf("NewTicketDomain failed: %v", err)
		}
		if ticketDomain == nil {
			t.Error("NewTicketDomain returned nil ticketDomain")
		}
	})

	t.Run("nil database", func(t *testing.T) {
		ticketDomain, err := main.NewTicketDomain(nil, bucket, privateKey, publicKey, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil database")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil database")
		}
	})

	t.Run("nil bucket", func(t *testing.T) {
		ticketDomain, err := main.NewTicketDomain(db, nil, privateKey, publicKey, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil bucket")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil bucket")
		}
	})

	t.Run("nil private key", func(t *testing.T) {
		ticketDomain, err := main.NewTicketDomain(db, bucket, nil, publicKey, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil private key")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil private key")
		}
	})

	t.Run("nil public key", func(t *testing.T) {
		ticketDomain, err := main.NewTicketDomain(db, bucket, privateKey, nil, mailer)
		if err == nil {
			t.Error("NewTicketDomain did not return error with nil public key")
		}
		if ticketDomain != nil {
			t.Error("NewTicketDomain returned non-nil ticketDomain with nil public key")
		}
	})

	t.Run("nil mailer", func(t *testing.T) {
		ticketDomain, err := main.NewTicketDomain(db, bucket, privateKey, publicKey, nil)
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

	ticketDomain, err := main.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender)
	if err != nil {
		t.Fatalf("creating a ticket domain instance: %s", err.Error())
	}

	userDomain := main.NewUserDomain(database)

	t.Run("Invalid Email and photo", func(t *testing.T) {
		err := ticketDomain.StorePaymentReceipt(context.Background(), "", nil, "")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		var validationError *main.ValidationError
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
		err := userDomain.CreateParticipant(ctx, main.CreateParticipantRequest{
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
		err := userDomain.CreateParticipant(ctx, main.CreateParticipantRequest{
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

		if err != nil && !errors.Is(err, main.ErrUserEmailNotFound) {
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

	ticketDomain, err := main.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender)
	if err != nil {
		t.Fatalf("creating a ticket domain instance: %s", err.Error())
	}

	userDomain := main.NewUserDomain(database)

	t.Run("Invalid email", func(t *testing.T) {
		_, err := ticketDomain.ValidatePaymentReceipt(context.Background(), "")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		var validationError *main.ValidationError
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

		if !errors.Is(err, main.ErrInvalidTicket) {
			t.Errorf("expecting an error of ErrInvalidTicket, instead got %s", err.Error())
		}
	})

	t.Run("Happy scenario", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		email := "johndoe+happy@example.com"
		err := userDomain.CreateParticipant(ctx, main.CreateParticipantRequest{
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
	ticketDomain, err := main.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender)
	if err != nil {
		t.Fatalf("creating a ticket domain instance: %s", err.Error())
	}
	t.Run("Invalid email", func(t *testing.T) {
		err := ticketDomain.VerifyIsStudent(context.Background(), "")
		if err == nil {
			t.Error("expecting an error, got nil instead")
		}

		var validationError *main.ValidationError
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
