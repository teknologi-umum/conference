package ticketing

import (
	"context"
	"fmt"
	"io"
	"mime"
	"time"

	"conf/user"
	"github.com/getsentry/sentry-go"
	"gocloud.dev/blob"
)

// StorePaymentReceipt stores the photo and email combination into our datastore.
// This will be reviewed manually by the TeknumConf team.
func (t *TicketDomain) StorePaymentReceipt(ctx context.Context, user user.User, photo io.Reader, contentType string) error {
	span := sentry.StartSpan(ctx, "ticketing.store_payment_receipt", sentry.WithTransactionName("StorePaymentReceipt"))
	defer span.Finish()

	var validationError ValidationError
	if photo == nil {
		validationError.Errors = append(validationError.Errors, "photo is nil")
	}

	if contentType == "" {
		validationError.Errors = append(validationError.Errors, "contentType is empty")
	}

	if len(validationError.Errors) > 0 {
		return validationError
	}

	// Write entry to database
	fileExtensions, _ := mime.ExtensionsByType(contentType)
	if len(fileExtensions) == 0 {
		fileExtensions = []string{""} // length is not zero, we can safely call fileExtensions[0]
	}

	// Store photo to filesystem (please use this one https://pkg.go.dev/gocloud.dev@v0.34.0/blob)
	blobKey := fmt.Sprintf("%s_%s.%s", time.Now().Format(time.RFC3339), user.Email, fileExtensions[0])
	err := t.bucket.Upload(ctx, blobKey, photo, &blob.WriterOptions{
		ContentType: contentType,
		Metadata: map[string]string{
			"email": user.Email,
		},
	})
	if err != nil {
		return fmt.Errorf("uploading to bucket storage: %w", err)
	}

	err = t.db.CreateTableRecords(ctx, t.tableId, []any{Ticketing{
		Email:            user.Email,
		ReceiptPhotoPath: blobKey,
		Paid:             false,
		SHA256Sum:        "",
		Used:             false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}})
	if err != nil {
		return fmt.Errorf("inserting ticketing entry into database: %w", err)
	}

	return nil
}
