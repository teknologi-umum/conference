package administrator

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"conf/administrator/jwt"
	"github.com/pquerna/otp/totp"
)

type Administrator struct {
	Username       string `yaml:"username"`
	HashedPassword string `yaml:"hashed_password"`
	TotpSecret     string `yaml:"totp_secret"`
}

func GenerateSecret(username string) (secret string, url string, err error) {
	generate, err := totp.Generate(totp.GenerateOpts{Issuer: "teknumconf", AccountName: username, Rand: rand.Reader})
	if err != nil {
		return "", "", err
	}

	return generate.Secret(), generate.URL(), nil
}

type AdministratorDomain struct {
	jwt            *jwt.JsonWebToken
	administrators []Administrator
}

func NewAdministratorDomain(administrators []Administrator) (*AdministratorDomain, error) {
	// Generate ed25519 key pairs for access and refresh tokens
	accessPublicKey, accessPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("generating fresh access key pair: %w", err)
	}

	refreshPublicKey, refreshPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("generating fresh refresh key pair: %w", err)
	}

	var randomIssuer = make([]byte, 18)
	_, _ = rand.Read(randomIssuer)

	var randomSubject = make([]byte, 16)
	_, _ = rand.Read(randomSubject)

	var randomAudience = make([]byte, 32)
	_, _ = rand.Read(randomAudience)

	authJwt := jwt.NewJwt(accessPrivateKey, accessPublicKey, refreshPrivateKey, refreshPublicKey, hex.EncodeToString(randomIssuer), hex.EncodeToString(randomSubject), hex.EncodeToString(randomAudience))

	return &AdministratorDomain{
		jwt:            authJwt,
		administrators: administrators,
	}, nil
}
