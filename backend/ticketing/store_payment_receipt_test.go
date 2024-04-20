package ticketing_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"strings"
	"testing"
	"time"

	"conf/ticketing"
	"conf/user"
)

func TestTicketDomain_StorePaymentReceipt(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating new ed25519 key: %s", err.Error())
		return
	}

	ticketDomain, err := ticketing.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender, tableId)
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
