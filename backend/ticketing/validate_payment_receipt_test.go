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

func TestTicketDomain_ValidatePaymentReceipt(t *testing.T) {
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

	t.Run("Email not found", func(t *testing.T) {
		t.Skip("We don't do any parameter checking on mock server")

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
