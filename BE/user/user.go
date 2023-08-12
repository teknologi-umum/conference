package user

import (
	"errors"

	"github.com/jackc/pgx"
)

type User struct {
	DB *pgx.Conn
}

func New(DB *pgx.Conn) *User {
	return &User{DB: DB}
}

type UserType string

var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidEmail    = errors.New("invalid email")

	UserTypeParticipant UserType = "participant"
	UserTypeSpeaker     UserType = "speaker"
)

type CreateParticipantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c CreateParticipantRequest) validate() error {
	if c.Name == "" {
		return ErrInvalidUserName
	}
	if c.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

func (u *User) CreateParticipant(req CreateParticipantRequest) error {
	if err := req.validate(); err != nil {
		return err
	}
	_, err := u.DB.Exec("INSERT INTO users (name, email, type) VALUES ($1, $2, $3)", req.Name, req.Email, UserTypeParticipant)
	if err != nil {
		return err
	}
	return nil
}
