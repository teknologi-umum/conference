package user

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDomain struct {
	DB *pgxpool.Pool
}

func New(DB *pgxpool.Pool) *UserDomain {
	// TODO: should return error on nil DB, must guarantee returned User is not nil
	return &UserDomain{DB: DB}
}

type Type string

const (
	TypeParticipant Type = "participant"
	TypeSpeaker     Type = "speaker"
)

var (
	ErrValidation = errors.New("validation")
)

type CreateParticipant struct {
	Name  string
	Email string
}

type User struct {
	Name        string
	Email       string
	Type        Type
	IsProcessed bool
}

func (c CreateParticipant) validate() error {
	if c.Name == "" {
		return fmt.Errorf("%w: invalid name", ErrValidation)
	}
	if c.Email == "" {
		return fmt.Errorf("%w: invalid email", ErrValidation)
	}
	return nil
}

func (u *UserDomain) CreateParticipant(ctx context.Context, req CreateParticipant) error {
	if err := req.validate(); err != nil {
		return err
	}

	c, err := u.DB.Acquire(ctx)
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
			return fmt.Errorf("%w (%s)", e, err.Error())
		}
		return err
	}

	err = t.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type UserFilter struct {
	Type        Type
	IsProcessed bool
}

func (u *UserDomain) GetUsers(ctx context.Context, filter UserFilter) ([]User, error) {
	c, err := u.DB.Acquire(ctx)
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
		err = rows.Scan(&user.Name, &user.Email, &user.Type, &user.IsProcessed)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserDomain) ExportUnprocessedUsersAsCSV(ctx context.Context) error {
	users, err := u.GetUsers(ctx, UserFilter{
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
			fmt.Sprintf("%t", user.IsProcessed),
		})
	}

	csvFile, err := os.Create("users.csv")
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
