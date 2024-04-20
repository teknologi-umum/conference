package administrator

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func (a *AdministratorDomain) Authenticate(ctx context.Context, username string, plainPassword string, otpCode string) (string, bool, error) {
	span := sentry.StartSpan(ctx, "administrator.authenticate", sentry.WithTransactionName("Authenticate"))
	defer span.Finish()

	var administrator Administrator
	for _, adm := range a.administrators {
		if adm.Username == username {
			administrator = adm
			break
		}
	}

	if administrator.Username == "" {
		return "", false, nil
	}

	hashedPassword, err := hex.DecodeString(administrator.HashedPassword)
	if err != nil {
		return "", false, fmt.Errorf("invalid hex string")
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", false, nil
		}

		return "", false, fmt.Errorf("password: %w", err)
	}

	ok := totp.Validate(otpCode, administrator.TotpSecret)
	if !ok {
		return "", false, nil
	}

	token, err := a.jwt.Sign(username)
	if err != nil {
		return "", false, fmt.Errorf("signing token: %w", err)
	}

	return token, true, nil
}
