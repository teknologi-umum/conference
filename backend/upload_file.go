package main

import (
	"crypto/ed25519"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UploadDomain struct {
	db         *pgxpool.Pool
	mailer     *Mailer
	privateKey *ed25519.PrivateKey
	publicKey  *ed25519.PublicKey
}

func NewUploadDomain(db *pgxpool.Pool, mailer *Mailer, privateKey *ed25519.PrivateKey, publicKey *ed25519.PublicKey) *UploadDomain {
	if db == nil {
		panic("db is nil")
	}
	if mailer == nil {
		panic("mailer is nil")
	}
	if privateKey == nil {
		panic("privateKey is nil")
	}
	if publicKey == nil {
		panic("publicKey is nil")
	}
	return &UploadDomain{db: db, mailer: mailer, privateKey: privateKey, publicKey: publicKey}
}

func (u *UploadDomain) UploadFile() {
	panic("TODO: implement me")
}

func (u *UploadDomain) ValidateFile() {
	panic("TODO: implement me")
}

func (u *UploadDomain) VerifyTransfer() {
	panic("TODO: implement me")
}
