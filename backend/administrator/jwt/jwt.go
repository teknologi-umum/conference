package jwt

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JsonWebToken struct {
	accessPrivateKey  ed25519.PrivateKey
	accessPublicKey   ed25519.PublicKey
	refreshPrivateKey ed25519.PrivateKey
	refreshPublicKey  ed25519.PublicKey
	issuer            string
	subject           string
	audience          string
}

func NewJwt(accessPrivateKey []byte, accessPublicKey []byte, refreshPrivateKey []byte, refreshPublicKey []byte, issuer string, subject string, audience string) *JsonWebToken {
	return &JsonWebToken{
		accessPrivateKey:  accessPrivateKey,
		accessPublicKey:   accessPublicKey,
		refreshPrivateKey: refreshPrivateKey,
		refreshPublicKey:  refreshPublicKey,
		issuer:            issuer,
		subject:           subject,
		audience:          audience,
	}
}

func (j *JsonWebToken) Sign(userId string) (accessToken string, err error) {
	accessRandId := make([]byte, 32)
	_, _ = rand.Read(accessRandId)

	accessClaims := jwt.MapClaims{
		"iss": j.issuer,
		"sub": j.subject,
		"aud": j.audience,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"jti": string(accessRandId),
		"uid": userId,
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodEdDSA, accessClaims).SignedString(j.accessPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return accessToken, nil
}

var ErrInvalidSigningMethod = errors.New("invalid signing method")
var ErrExpired = errors.New("token expired")
var ErrInvalid = errors.New("token invalid")
var ErrClaims = errors.New("token claims invalid")

func (j *JsonWebToken) VerifyAccessToken(token string) (userId string, err error) {
	if token == "" {
		return "", ErrInvalid
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodEd25519)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}
		return j.accessPublicKey, nil
	})
	if err != nil {
		if parsedToken != nil && !parsedToken.Valid {
			// Check if the error is a type of jwt.ValidationError
			validationError, ok := err.(*jwt.ValidationError)
			if ok {
				if validationError.Errors&jwt.ValidationErrorExpired != 0 {
					return "", ErrExpired
				}

				if validationError.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					return "", ErrInvalid
				}

				if validationError.Errors&jwt.ValidationErrorClaimsInvalid != 0 {
					return "", ErrClaims
				}

				return "", fmt.Errorf("failed to parse access token: %w", err)
			}

			return "", fmt.Errorf("non-validation error during parsing token: %w", err)
		}

		return "", fmt.Errorf("token is valid or parsedToken is not nil: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrClaims
	}

	if !claims.VerifyAudience(j.audience, true) {
		return "", ErrInvalid
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", ErrExpired
	}

	if !claims.VerifyIssuer(j.issuer, true) {
		return "", ErrInvalid
	}

	if !claims.VerifyNotBefore(time.Now().Unix(), true) {
		return "", ErrInvalid
	}

	jwtId, ok := claims["jti"].(string)
	if !ok {
		return "", ErrClaims
	}

	if jwtId == "" {
		return "", ErrClaims
	}

	userId, ok = claims["uid"].(string)
	if !ok {
		return "", ErrClaims
	}

	return userId, nil
}
