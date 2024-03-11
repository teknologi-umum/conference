package administrator

import (
	"context"

	"github.com/getsentry/sentry-go"
)

func (a *AdministratorDomain) Validate(ctx context.Context, token string) (Administrator, bool, error) {
	span := sentry.StartSpan(ctx, "administrator.validate", sentry.WithTransactionName("Validate"))
	defer span.Finish()

	if token == "" {
		return Administrator{}, false, nil
	}

	username, err := a.jwt.VerifyAccessToken(token)
	if err != nil {
		return Administrator{}, false, nil
	}

	var administrator Administrator
	for _, adm := range a.administrators {
		if adm.Username == username {
			administrator = adm
			break
		}
	}

	return administrator, true, nil
}
