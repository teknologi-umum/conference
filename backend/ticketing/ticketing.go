package ticketing

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"strings"
	"time"

	"conf/mailer"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skip2/go-qrcode"
	"gocloud.dev/blob"
)

type TicketDomain struct {
	db         *pgxpool.Pool
	bucket     *blob.Bucket
	privateKey *ed25519.PrivateKey
	publicKey  *ed25519.PublicKey
	mailer     *mailer.Mailer
}

func NewTicketDomain(db *pgxpool.Pool, bucket *blob.Bucket, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, mailer *mailer.Mailer) (*TicketDomain, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	if bucket == nil {
		return nil, fmt.Errorf("bucket is nil")
	}

	if privateKey == nil {
		return nil, fmt.Errorf("privateKey is nil")
	}

	if publicKey == nil {
		return nil, fmt.Errorf("publicKey is nil")
	}

	if mailer == nil {
		return nil, fmt.Errorf("mailer is nil")
	}

	return &TicketDomain{
		db:         db,
		bucket:     bucket,
		privateKey: &privateKey,
		publicKey:  &publicKey,
		mailer:     mailer,
	}, nil
}

// StorePaymentReceipt stores the photo and email combination into our datastore.
// This will be reviewed manually by the TeknumConf team.
func (t *TicketDomain) StorePaymentReceipt(ctx context.Context, email string, photo io.Reader, contentType string) error {
	span := sentry.StartSpan(ctx, "ticket.store_payment_receipt")
	defer span.Finish()

	var validationError ValidationError
	if email == "" {
		validationError.Errors = append(validationError.Errors, "email is empty")
	}

	if photo == nil {
		validationError.Errors = append(validationError.Errors, "photo is nil")
	}

	if contentType == "" {
		validationError.Errors = append(validationError.Errors, "contentType is empty")
	}

	if len(validationError.Errors) > 0 {
		return validationError
	}

	// Write entry to postgres
	conn, err := t.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	var userID string
	err = tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if e := tx.Rollback(ctx); e != nil {
				return fmt.Errorf("rolling back transaction: %w (%s)", e, ErrUserEmailNotFound.Error())
			}

			return ErrUserEmailNotFound
		}

		if e := tx.Rollback(ctx); e != nil {
			return fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return fmt.Errorf("executing select query: %w", err)
	}

	fileExtensions, _ := mime.ExtensionsByType(contentType)
	if len(fileExtensions) == 0 {
		fileExtensions = []string{""} // length is not zero, we can safely call fileExtensions[0]
	}

	// Store photo to filesystem (please use this one https://pkg.go.dev/gocloud.dev@v0.34.0/blob)
	blobKey := fmt.Sprintf("%s_%s.%s", time.Now().Format(time.RFC3339), email, fileExtensions[0])
	err = t.bucket.Upload(ctx, blobKey, photo, &blob.WriterOptions{
		ContentType: contentType,
		Metadata: map[string]string{
			"email": email,
		},
	})
	if err != nil {
		return fmt.Errorf("uploading to bucket storage: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO ticketing (id, email, receipt_photo_path) VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET receipt_photo_path = $3, updated_at = NOW()`,
		uuid.New(),
		email,
		blobKey)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return fmt.Errorf("executing query: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

// ValidatePaymentReceipt marks an email payment status as paid. It will create a signature using Ed25519,
// encode it to a QRCode image, and send the QRCode to the user's email. It returns hex-encoded SHA256SUM
// of the QR code.
//
// It will return ErrInvalidTicket if the payment receipt's not uploaded yet.
func (t *TicketDomain) ValidatePaymentReceipt(ctx context.Context, email string) (string, error) {
	span := sentry.StartSpan(ctx, "ticket.validate_payment_receipt")
	defer span.Finish()

	if email == "" {
		return "", ValidationError{Errors: []string{"email is empty"}}
	}

	// Mark payment status as paid on postgres
	conn, err := t.db.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return "", fmt.Errorf("creating transaction: %w", err)
	}

	var id uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM ticketing WHERE email = $1`, email).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w: not exists", ErrInvalidTicket)
		}

		if e := tx.Rollback(ctx); e != nil {
			return "", fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return "", fmt.Errorf("executing select query: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE ticketing SET paid = TRUE, updated_at = NOW() WHERE email = $1`, email)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return "", fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return "", fmt.Errorf("executing update query: %w", err)
	}

	// Create a signature using unique key based on the email and random id combination (possibly using any non-text based encoding)
	sha384Hasher := sha512.New384()
	sha384Hasher.Write([]byte(email))
	hashedEmail := sha384Hasher.Sum(nil)
	payload := fmt.Sprintf("%s:%s", id.String(), base64.StdEncoding.EncodeToString(hashedEmail))

	signature := ed25519.Sign(*t.privateKey, []byte(payload))

	// Generate QR code with https://github.com/skip2/go-qrcode
	qrImage, err := qrcode.Encode(fmt.Sprintf("%s;%s", hex.EncodeToString(signature), payload), qrcode.High, 1024)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return "", fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return "", fmt.Errorf("generating qr code: %w", err)
	}

	// Create SHA256SUM to the generated QR code
	sha256Hasher := sha256.New()
	sha256Hasher.Write(qrImage)
	sha256Sum := sha256Hasher.Sum(nil)

	imageCid, _, _ := strings.Cut(uuid.NewString(), "-")

	// Send email programmatically
	err = t.mailer.Send(ctx, &mailer.Mail{
		RecipientName:  "",
		RecipientEmail: email,
		Subject:        "TeknumConf 2023: Tiket Anda!",
		PlainTextBody: `Hai! Ini dia email yang kamu tunggu-tunggu💃
        
Pembayaran kamu telah di konfirmasi! Dibawah ini terdapat QR code sebagai tiket kamu masuk ke TeknumConf 2023.
Apabila kamu mendapat student discount, pastikan kamu membawa Kartu Mahasiswa atau Kartu Pelajar ya!
Panitia akan melakukan verifikasi tambahan pada lokasi untuk memastikan kalau kamu betulan pelajar.

Sampai jumpa di TeknumConf 2023!

Email ini hanya tertuju untuk Anda. Apabila Anda merasa tidak mendaftar untuk TeknumConf 2023,
harap abaikan email ini. Terima kasih!`,
		HtmlBody: `<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml">
    <head>
        <meta content="IE=edge" http-equiv="X-UA-Compatible" />
        <meta content="width=device-width,initial-scale=1 user-scalable=yes" name="viewport" />
        <meta content="telephone=no, date=no, address=no, email=no, url=no" name="format-detection" />
        <meta name="x-apple-disable-message-reformatting" />
        <meta charset="UTF-8" />
        <!--[if mso]>
            <noscript>
                <xml>
                    <o:OfficeDocumentSettings>
                        <o:PixelsPerInch>96</o:PixelsPerInch>
                    </o:OfficeDocumentSettings>
                </xml>
            </noscript>
        <![endif]-->

        <style>
            * {
                font-family: 'Rubik', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            }
        </style>
        <title>TeknumConf 2023: Tiket Anda!</title>
    </head>
    <body>
        <h1>Hai! Ini dia email yang kamu tunggu-tunggu💃</h1>
        <p>
            Pembayaran kamu telah di konfirmasi! Dibawah ini terdapat QR code sebagai tiket kamu masuk ke TeknumConf 2023.
            Apabila kamu mendapat <i>student discount</i>, pastikan kamu membawa Kartu Mahasiswa atau Kartu Pelajar ya!
            Panitia akan melakukan verifikasi tambahan pada lokasi untuk memastikan kalau kamu betulan pelajar.
        </p>
        <p><b>Sampai jumpa di TeknumConf 2023!</b></p>
        <p><img src="cid:` + imageCid + `" style="width: 100%; max-width: 720px;"></p>
        <p>
            <small>
                Email ini hanya tertuju untuk Anda. Apabila Anda merasa tidak mendaftar untuk TeknumConf 2023,
                harap abaikan email ini. Terima kasih!
            </small>
        </p>
    </body>
</html>
`,
		Attachments: []mailer.Attachment{
			{
				Name:               "qrcode_ticket.png",
				Description:        "QR code ticket TeknumConf 2023",
				ContentType:        "image/png",
				ContentDisposition: mailer.ContentDispositionInline,
				ContentId:          imageCid,
				SHA256Checksum:     sha256Sum,
				Payload:            qrImage,
			},
		},
	})
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return "", fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return "", fmt.Errorf("sending mail: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE ticketing SET sha256sum = $1, updated_at = NOW() WHERE email = $2`, sha256Sum, email)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return "", fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return "", fmt.Errorf("executing select query: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("commiting transaction: %w", err)
	}

	return hex.EncodeToString(sha256Sum), nil
}

// VerifyTicket will verify a ticket from the QR code payload. It will disassemble the payload and validate
// the signature and mark the ticket as used. Each ticket can only be used once.
//
// If the signature is invalid or the ticket is used, it will return ErrInvalidTicket error.
func (t *TicketDomain) VerifyTicket(ctx context.Context, payload []byte) (email string, name string, student bool, err error) {
	span := sentry.StartSpan(ctx, "ticket.verify_ticket")
	defer span.Finish()

	if len(payload) == 0 {
		return "", "", false, ValidationError{Errors: []string{"payload is empty"}}
	}

	// Separate the payload into the signature + email + random id that's generated from ValidatePaymentReceipt
	rawSignature, payloadAfter, found := bytes.Cut(payload, []byte(";"))
	if !found {
		return "", "", false, ErrInvalidTicket
	}

	rawTicketId, rawHashedEmail, found := bytes.Cut(payloadAfter, []byte(":"))
	if !found {
		return "", "", false, ErrInvalidTicket
	}

	ticketId, err := uuid.ParseBytes(rawTicketId)
	if err != nil {
		return "", "", false, ErrInvalidTicket
	}

	userHashedEmail, err := base64.StdEncoding.DecodeString(string(rawHashedEmail))
	if err != nil {
		return "", "", false, fmt.Errorf("decoding base64 string for email: %w", err)
	}

	signature, err := hex.DecodeString(string(rawSignature))
	if err != nil {
		return "", "", false, fmt.Errorf("decoding hex string for signature: %w", err)
	}

	// Validate the signature and its message using ed25519. If it's invalid, return ErrInvalidTicket
	signatureValidated := ed25519.Verify(*t.publicKey, payloadAfter, signature)
	if !signatureValidated {
		return "", "", false, fmt.Errorf("%w (verifying signature)", ErrInvalidTicket)
	}

	// Check the ticket if it's been used before. If it is, return ErrInvalidTicket. Decorate it a bit.
	conn, err := t.db.Acquire(ctx)
	if err != nil {
		return "", "", false, fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return "", "", false, fmt.Errorf("creating transaction: %w", err)
	}

	err = tx.QueryRow(ctx, "SELECT email, student FROM ticketing WHERE id = $1 AND used = FALSE", ticketId).Scan(&email, &student)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return "", "", false, fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", false, fmt.Errorf("%w (id not exists, or ticket has been used)", ErrInvalidTicket)
		}

		return "", "", false, fmt.Errorf("acquiring data from table: %w", err)
	}

	err = tx.QueryRow(ctx, "SELECT name FROM users WHERE email = $1", email).Scan(&name)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return "", "", false, fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		// Do not return error if something's wrong with name.
		// Instead, just report to Sentry.
		if !errors.Is(err, pgx.ErrNoRows) {
			sentry.GetHubFromContext(ctx).CaptureException(err)
		}
	}

	// Validate email
	sha384Hasher := sha512.New384()
	sha384Hasher.Write([]byte(email))
	hashedEmail := sha384Hasher.Sum(nil)
	if !bytes.Equal(hashedEmail, userHashedEmail) {
		return "", "", false, fmt.Errorf("%w (mismatched email)", ErrInvalidTicket)
	}

	// Mark the ticket as used
	_, err = tx.Exec(ctx, "UPDATE ticketing SET used = TRUE WHERE id = $1", ticketId)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return "", "", false, fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return "", "", false, fmt.Errorf("acquiring data from table: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", false, fmt.Errorf("commiting transaction: %w", err)
	}

	return email, name, student, nil
}

func (t *TicketDomain) VerifyIsStudent(ctx context.Context, email string) (err error) {
	span := sentry.StartSpan(ctx, "ticket.verify_is_student")
	defer span.Finish()

	if email == "" {
		return ValidationError{Errors: []string{"email is empty"}}
	}

	c, err := t.db.Acquire(ctx)
	if err != nil {
		return
	}
	defer c.Release()

	tx, err := c.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "UPDATE ticketing SET student = TRUE WHERE email = $1", email)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			err = e
			return
		}

		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}

	return nil
}
