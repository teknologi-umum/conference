package main

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDomain struct {
	db *pgxpool.Pool
}

func NewUserDomain(db *pgxpool.Pool) *UserDomain {
	if db == nil {
		panic("db is nil")
	}

	return &UserDomain{db: db}
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

	c, err := u.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	t, err := c.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}

	_, err = t.Exec(
		ctx,
		"INSERT INTO users (name, email, type) VALUES ($1, $2, $3)",
		req.Name,
		req.Email,
		TypeParticipant,
	)
	if err != nil {
		if e := t.Rollback(ctx); e != nil {
			return e
		}
		return err
	}

	err = t.Commit(ctx)
	if err != nil {
		return err
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

	c, err := u.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()

	rows, err := c.Query(
		ctx,
		"SELECT name, email, type, is_processed FROM users WHERE type = $1 AND is_processed = $2",
		filter.Type,
		filter.IsProcessed,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Name, &user.Email, &user.Type, &user.IsProcessed)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
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
	defer csvFile.Close()

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
