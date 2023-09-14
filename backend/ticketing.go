package main

import (
	"context"
	"crypto/ed25519"
	"errors"
	"io"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidTicket = errors.New("invalid ticket")

type TicketDomain struct {
	db         *pgxpool.Pool
	privateKey *ed25519.PrivateKey
	publicKey  *ed25519.PublicKey
	mailer     *Mailer
}

func NewTicketDomain(db *pgxpool.Pool) (*TicketDomain, error) {
	return &TicketDomain{
		db:         nil, // TODO: fill these
		privateKey: nil,
		publicKey:  nil,
	}, nil
}

// StorePaymentReceipt stores the photo and email combination into our datastore.
// This will be reviewed manually by the TeknumConf team.
func (t *TicketDomain) StorePaymentReceipt(ctx context.Context, email string, photo io.ReadCloser) error {
	// Write entry to postgres
	// Store photo to filesystem (please use this one https://pkg.go.dev/gocloud.dev@v0.34.0/blob)
	panic("TODO: implement me")
}

// ValidatePaymentReceipt marks an email payment status as paid. It will create a signature using Ed25519,
// encode it to a QRCode image, and send the QRCode to the user's email.
func (t *TicketDomain) ValidatePaymentReceipt(ctx context.Context, email string) (sha256sum string, err error) {
	// Mark payment status as paid on postgres
	// Create a signature using unique key based on the email and random id combination (possibly using any non-text based encoding)
	// Generate QR code with https://github.com/skip2/go-qrcode
	// Send email programmatically
	// Create SHA256SUM to the generated QR code
	panic("TODO: implement me")
}

// VerifyTicket will verify a ticket from the QR code payload. It will disassemble the payload and validate
// the signature and mark the ticket as used. Each ticket can only be used once.
//
// If the signature is invalid or the ticket is used, it will return ErrInvalidTicket error.
func (t *TicketDomain) VerifyTicket(ctx context.Context, payload []byte) (ok bool, err error) {
	// Separate the payload into the signature + email + random id that's generated from ValidatePaymentReceipt
	// Validate the signature and its message using ed25519. If it's invalid, return ErrInvalidTicket
	// Check the ticket if it's been used before. If it is, return ErrInvalidTicket. Decorate it a bit.
	// Mark the ticket as used
	panic("TODO: implement me")
}
