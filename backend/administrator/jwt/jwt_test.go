package jwt_test

import (
	"crypto/ed25519"
	"errors"
	"log"
	"os"
	"testing"

	"conf/administrator/jwt"
)

var authJwt *jwt.JsonWebToken

func TestMain(m *testing.M) {
	// Generate ed25519 key pairs for access and refresh tokens
	accessPublicKey, accessPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("failed to generate access key pair: %v", err)
	}

	refreshPublicKey, refreshPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("failed to generate refresh key pair: %v", err)
	}

	authJwt = jwt.NewJwt(accessPrivateKey, accessPublicKey, refreshPrivateKey, refreshPublicKey, "kodiiing", "user", "kodiiing")

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestSign(t *testing.T) {
	accessToken, err := authJwt.Sign("john")
	if err != nil {
		t.Errorf("failed to sign access token: %v", err)
	}

	if accessToken == "" {
		t.Error("access token is empty")
	}
}

func TestVerify(t *testing.T) {
	accessToken, err := authJwt.Sign("john")
	if err != nil {
		t.Errorf("failed to sign access token: %v", err)
	}

	if accessToken == "" {
		t.Error("access token is empty")
	}

	accessId, err := authJwt.VerifyAccessToken(accessToken)
	if err != nil {
		t.Errorf("failed to verify access token: %v", err)
	}

	if accessId != "john" {
		t.Errorf("access id is not 'john': %v", accessId)
	}
}

func TestVerifyEmpty(t *testing.T) {
	accessId, err := authJwt.VerifyAccessToken("")
	if err == nil {
		t.Errorf("access token is valid: %v", accessId)
	}

	if !errors.Is(err, jwt.ErrInvalid) {
		t.Errorf("error is not ErrInvalid: %v", err)
	}
}
