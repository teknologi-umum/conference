package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	DB *pgxpool.Pool
}

func New(DB *pgxpool.Pool) *User {
	// TODO: should return error on nil DB, must guarantee returned User is not nil
	return &User{DB: DB}
}

type Type string

const (
	TypeParticipant Type = "participant"
	TypeSpeaker     Type = "speaker"
)

var (
	ErrValidation = errors.New("validation")
)

type CreateParticipantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c CreateParticipantRequest) validate() error {
	if c.Name == "" {
		return fmt.Errorf("%w: invalid name", ErrValidation)
	}
	if c.Email == "" {
		return fmt.Errorf("%w: invalid email", ErrValidation)
	}
	return nil
}

func (u *User) CreateParticipant(ctx context.Context, req CreateParticipantRequest) error {
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
