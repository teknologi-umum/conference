package ticketing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"conf/nocodb"
	"conf/user"
	"github.com/getsentry/sentry-go"
)

func (t *TicketDomain) VerifyIsStudent(ctx context.Context, user user.User) (err error) {
	span := sentry.StartSpan(ctx, "ticketing.verify_is_student", sentry.WithTransactionName("VerifyIsStudent"))
	defer span.Finish()

	var rawTicketingResults []Ticketing
	_, err = t.db.ListTableRecords(ctx, "TODO: Table Id", &rawTicketingResults, nocodb.ListTableRecordOptions{
		Where: fmt.Sprintf("(Email,eq,%s)", user.Email),
		Sort:  []nocodb.Sort{nocodb.SortDescending("CreatedAt")},
		Limit: 1,
	})
	if err != nil {
		return fmt.Errorf("acquiring records: %w", err)
	}

	if len(rawTicketingResults) == 0 {
		return fmt.Errorf("%w: not exists", ErrInvalidTicket)
	}

	var ticketing = rawTicketingResults[0]

	err = t.db.UpdateTableRecords(ctx, "TODO: Table Id", []any{NullTicketing{
		Id:        sql.NullInt64{Int64: ticketing.Id, Valid: true},
		Student:   sql.NullBool{Bool: true, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}})
	if err != nil {
		return fmt.Errorf("updating table records: %w", err)
	}

	return nil
}
