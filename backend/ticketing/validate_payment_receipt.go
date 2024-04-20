package ticketing

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"conf/mailer"
	"conf/nocodb"
	"conf/user"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// ValidatePaymentReceipt marks an email payment status as paid. It will create a signature using Ed25519,
// encode it to a QRCode image, and send the QRCode to the user's email. It returns hex-encoded SHA256SUM
// of the QR code.
//
// It will return ErrInvalidTicket if the payment receipt's not uploaded yet.
func (t *TicketDomain) ValidatePaymentReceipt(ctx context.Context, user user.User) (string, error) {
	span := sentry.StartSpan(ctx, "ticketing.validate_payment_receipt", sentry.WithTransactionName("ValidatePaymentReceipt"))
	defer span.Finish()

	// Mark payment status as paid on database
	var rawTicketingResults []Ticketing
	_, err := t.db.ListTableRecords(ctx, t.tableId, &rawTicketingResults, nocodb.ListTableRecordOptions{
		Where: fmt.Sprintf("(Email,eq,%s)", user.Email),
		Sort:  []nocodb.Sort{nocodb.SortDescending("CreatedAt")},
		Limit: 1,
	})
	if err != nil {
		return "", fmt.Errorf("acquiring records: %w", err)
	}

	if len(rawTicketingResults) == 0 {
		return "", fmt.Errorf("%w: not exists", ErrInvalidTicket)
	}

	var ticketing = rawTicketingResults[0]

	// Create a signature using unique key based on the email and random id combination (possibly using any non-text based encoding)
	sha384Hasher := sha512.New384()
	sha384Hasher.Write([]byte(user.Email))
	hashedEmail := sha384Hasher.Sum(nil)
	payload := fmt.Sprintf("%s:%s", strconv.FormatInt(ticketing.Id, 10), base64.StdEncoding.EncodeToString(hashedEmail))

	signature := ed25519.Sign(*t.privateKey, []byte(payload))

	// Generate QR code with https://github.com/skip2/go-qrcode
	qrImage, err := qrcode.Encode(fmt.Sprintf("%s;%s", hex.EncodeToString(signature), payload), qrcode.High, 1024)
	if err != nil {
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
		RecipientEmail: user.Email,
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
		return "", fmt.Errorf("sending mail: %w", err)
	}

	err = t.db.UpdateTableRecords(ctx, t.tableId, []any{NullTicketing{
		Id:        sql.NullInt64{Int64: ticketing.Id, Valid: true},
		Paid:      sql.NullBool{Bool: true, Valid: true},
		SHA256Sum: sql.NullString{String: hex.EncodeToString(sha256Sum), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}})
	if err != nil {
		return "", fmt.Errorf("updating table records: %w", err)
	}

	return hex.EncodeToString(sha256Sum), nil
}
