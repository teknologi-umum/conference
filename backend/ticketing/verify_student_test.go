package ticketing_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"conf/ticketing"
	"conf/user"
)

func TestTicketDomain_VerifyIsStudent(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating new ed25519 key: %s", err.Error())
		return
	}

	ticketDomain, err := ticketing.NewTicketDomain(database, bucket, privateKey, publicKey, mailSender, tableId)
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
