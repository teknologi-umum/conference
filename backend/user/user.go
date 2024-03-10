package user

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"conf/nocodb"
	"github.com/getsentry/sentry-go"
)

type UserDomain struct {
	db      *nocodb.Client
	tableId string
}

func NewUserDomain(db *nocodb.Client, tableId string) (*UserDomain, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	if tableId == "" {
		return nil, fmt.Errorf("tableId is empty")
	}

	return &UserDomain{db: db, tableId: tableId}, nil
}

type Type string

const (
	TypeParticipant Type = "participant"
	TypeSpeaker     Type = "speaker"
)

type CreateParticipantRequest struct {
	Name  string
	Email string
}

type User struct {
	Name        string
	Email       string
	Type        Type
	IsProcessed bool
	CreatedAt   time.Time
}

func (c CreateParticipantRequest) validate() (errors []string) {
	if c.Name == "" {
		errors = append(errors, "Invalid name")
	}

	if c.Email == "" {
		errors = append(errors, "Invalid email")
	}

	return errors
}

func (u *UserDomain) CreateParticipant(ctx context.Context, req CreateParticipantRequest) error {
	span := sentry.StartSpan(ctx, "user.create_participant")
	defer span.Finish()

	if errors := req.validate(); len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}

	user := User{
		Name:        req.Name,
		Email:       req.Email,
		Type:        TypeParticipant,
		IsProcessed: false,
		CreatedAt:   time.Now(),
	}

	err := u.db.CreateTableRecords(ctx, u.tableId, []any{user})
	if err != nil {
		return fmt.Errorf("creating table records: %w", err)
	}

	return nil
}

type UserFilterRequest struct {
	Type        Type
	IsProcessed bool
}

func (u *UserDomain) GetUsers(ctx context.Context, filter UserFilterRequest) ([]User, error) {
	span := sentry.StartSpan(ctx, "user.get_users")
	defer span.Finish()

	var users []User
	var offset int64
	for {
		var currentUserSets []User
		pageInfo, err := u.db.ListTableRecords(ctx, u.tableId, &currentUserSets, nocodb.ListTableRecordOptions{
			Offset: offset,
			Where:  fmt.Sprintf("(Type,eq,%s)~and(IsProcessed,eq,%s)", filter.Type, strconv.FormatBool(filter.IsProcessed)),
		})
		if err != nil {
			return users, fmt.Errorf("list table records: %w", err)
		}

		offset += int64(len(currentUserSets))
		users = append(users, currentUserSets...)

		if pageInfo.IsLastPage {
			break
		}
	}

	return users, nil
}

func (u *UserDomain) ExportUnprocessedUsersAsCSV(ctx context.Context) error {
	span := sentry.StartSpan(ctx, "user.export_unprocessed_users_as_csv")
	defer span.Finish()

	users, err := u.GetUsers(ctx, UserFilterRequest{
		Type:        TypeParticipant,
		IsProcessed: false,
	})
	if err != nil {
		return err
	}

	csvData := [][]string{
		{"name", "email", "type", "is_processed"},
	}

	for _, user := range users {
		csvData = append(csvData, []string{
			user.Name,
			user.Email,
			string(user.Type),
			strconv.FormatBool(user.IsProcessed),
		})
	}

	csvFile, err := os.Create("/app/csv/users.csv")
	if err != nil {
		return err
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			sentry.GetHubFromContext(ctx).Scope().SetLevel(sentry.LevelWarning)
			sentry.GetHubFromContext(ctx).CaptureException(err)
		}
	}(csvFile)

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for _, row := range csvData {
		err := csvWriter.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
