package ticketing

import (
	"crypto/ed25519"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"conf/mailer"
	"conf/nocodb"

	"gocloud.dev/blob"
)

type TicketDomain struct {
	db         *nocodb.Client
	bucket     *blob.Bucket
	privateKey *ed25519.PrivateKey
	publicKey  *ed25519.PublicKey
	mailer     *mailer.Mailer
}

func NewTicketDomain(db *nocodb.Client, bucket *blob.Bucket, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, mailer *mailer.Mailer) (*TicketDomain, error) {
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

type Ticketing struct {
	Id               int64     `json:"Id,omitempty"`
	Email            string    `json:"Email,omitempty"`
	ReceiptPhotoPath string    `json:"ReceiptPhotoPath,omitempty"`
	Paid             bool      `json:"Paid,omitempty"`
	Student          bool      `json:"Student,omitempty"`
	SHA256Sum        string    `json:"SHA256Sum,omitempty"`
	Used             bool      `json:"Used,omitempty"`
	CreatedAt        time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt        time.Time `json:"UpdatedAt,omitempty"`
}

type NullTicketing struct {
	Id               sql.NullInt64  `json:"Id,omitempty"`
	Email            sql.NullString `json:"Email,omitempty"`
	ReceiptPhotoPath sql.NullString `json:"ReceiptPhotoPath,omitempty"`
	Paid             sql.NullBool   `json:"Paid,omitempty"`
	Student          sql.NullBool   `json:"Student,omitempty"`
	SHA256Sum        sql.NullString `json:"SHA256Sum,omitempty"`
	Used             sql.NullBool   `json:"Used,omitempty"`
	CreatedAt        sql.NullTime   `json:"CreatedAt,omitempty"`
	UpdatedAt        sql.NullTime   `json:"UpdatedAt,omitempty"`
}

func (t NullTicketing) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	ut := reflect.TypeOf(t)
	uv := reflect.ValueOf(t)

	for i := 0; i < ut.NumField(); i++ {
		field := ut.Field(i)
		intf := uv.Field(i).Interface()
		valuer, ok := intf.(driver.Valuer)
		if ok {
			v, err := valuer.Value()
			if err != nil {
				continue
			}

			if v == nil {
				continue
			}

			m[field.Name] = v
			continue
		}

		m[field.Name] = intf
	}

	return json.Marshal(m)
}
