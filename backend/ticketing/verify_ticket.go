package ticketing

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"conf/nocodb"
	"github.com/getsentry/sentry-go"
)

// VerifyTicket will verify a ticket from the QR code payload. It will disassemble the payload and validate
// the signature and mark the ticket as used. Each ticket can only be used once.
//
// If the signature is invalid or the ticket is used, it will return ErrInvalidTicket error.
func (t *TicketDomain) VerifyTicket(ctx context.Context, payload []byte) (ticketing Ticketing, err error) {
	span := sentry.StartSpan(ctx, "ticketing.verify_ticket", sentry.WithTransactionName("VerifyTicket"))
	defer span.Finish()

	if len(payload) == 0 {
		return Ticketing{}, ValidationError{Errors: []string{"payload is empty"}}
	}

	// Separate the payload into the signature + email + random id that's generated from ValidatePaymentReceipt
	rawSignature, payloadAfter, found := bytes.Cut(payload, []byte(";"))
	if !found {
		return Ticketing{}, ErrInvalidTicket
	}

	rawTicketId, rawHashedEmail, found := bytes.Cut(payloadAfter, []byte(":"))
	if !found {
		return Ticketing{}, ErrInvalidTicket
	}

	ticketId, err := strconv.ParseInt(string(rawTicketId), 10, 64)
	if err != nil {
		return Ticketing{}, ErrInvalidTicket
	}

	userHashedEmail, err := base64.StdEncoding.DecodeString(string(rawHashedEmail))
	if err != nil {
		return Ticketing{}, fmt.Errorf("decoding base64 string for email: %w", err)
	}

	signature, err := hex.DecodeString(string(rawSignature))
	if err != nil {
		return Ticketing{}, fmt.Errorf("decoding hex string for signature: %w", err)
	}

	// Validate the signature and its message using ed25519. If it's invalid, return ErrInvalidTicket
	signatureValidated := ed25519.Verify(*t.publicKey, payloadAfter, signature)
	if !signatureValidated {
		return Ticketing{}, fmt.Errorf("%w (verifying signature)", ErrInvalidTicket)
	}

	// Check the ticket if it's been used before. If it is, return ErrInvalidTicket. Decorate it a bit.
	var rawTicketingResults []Ticketing
	_, err = t.db.ListTableRecords(ctx, t.tableId, &rawTicketingResults, nocodb.ListTableRecordOptions{
		Where: fmt.Sprintf("(Id,eq,%d)~and(Used,eq,false)", ticketId),
		Sort:  []nocodb.Sort{nocodb.SortDescending("CreatedAt")},
		Limit: 1,
	})
	if err != nil {
		return Ticketing{}, fmt.Errorf("acquiring records: %w", err)
	}

	if len(rawTicketingResults) == 0 {
		return Ticketing{}, fmt.Errorf("%w: not exists", ErrInvalidTicket)
	}

	ticketing = rawTicketingResults[0]

	// Validate email
	sha384Hasher := sha512.New384()
	sha384Hasher.Write([]byte(ticketing.Email))
	hashedEmail := sha384Hasher.Sum(nil)
	if !bytes.Equal(hashedEmail, userHashedEmail) {
		return Ticketing{}, fmt.Errorf("%w (mismatched email)", ErrInvalidTicket)
	}

	// Mark the ticket as used
	err = t.db.UpdateTableRecords(ctx, t.tableId, []any{NullTicketing{
		Id:        sql.NullInt64{Int64: ticketing.Id, Valid: true},
		Used:      sql.NullBool{Bool: true, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}})
	if err != nil {
		return Ticketing{}, fmt.Errorf("updating table records: %w", err)
	}

	return ticketing, nil
}
